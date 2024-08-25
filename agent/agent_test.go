package agent_test

import (
	"eapm-ebpf/agent/conn"
	"eapm-ebpf/bpf"
	"eapm-ebpf/common"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

func tearDown() {
	// Do something here.
	time.Sleep(1 * time.Second)
	fmt.Println("teardown!")
}

func TestConnectSyscall(t *testing.T) {
	connEventList := make([]bpf.AgentConnEvtT, 0)
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
		},
		&connEventList,
		nil,
		nil,
		nil,
		agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()
	fmt.Println("Start Send Http Request")
	sendTestRequest(t, SendTestHttpRequestOptions{disableKeepAlived: true})

	time.Sleep(1 * time.Second)
	if len(connEventList) == 0 {
		t.Fatalf("no conn event!")
	}
	connectEvent := findInterestedConnEvent(t, connEventList, FindInterestedConnEventOptions{
		findByRemotePort: true,
		findByConnType:   true,
		remotePort:       80,
		connType:         bpf.AgentConnTypeTKConnect,
		throw:            true})[0]
	AssertConnEvent(t, connectEvent, ConnEventAssertions{
		expectPid:             uint32(os.Getpid()),
		expectRemotePort:      80,
		expectLocalAddrFamily: common.AF_INET,
		expectRemoteFamily:    common.AF_INET,
		expectReadBytes:       0,
		expectWriteBytes:      0,
		expectLocalPort:       -1,
		expectConnEventType:   bpf.AgentConnTypeTKConnect,
	})

}

func TestCloseSyscall(t *testing.T) {
	connEventList := make([]bpf.AgentConnEvtT, 0)
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallCloseEntry,
			bpf.AttachSyscallCloseExit},
		&connEventList,
		nil,
		nil,
		nil, agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()
	fmt.Println("Start Send Http Request")
	sendTestRequest(t, SendTestHttpRequestOptions{disableKeepAlived: true})

	time.Sleep(1 * time.Second)
	if len(connEventList) == 0 {
		t.Fatalf("no conn event!")
	}
	connectEvent := findInterestedConnEvent(t, connEventList, FindInterestedConnEventOptions{
		findByRemotePort: true,
		findByConnType:   true,
		remotePort:       80,
		connType:         bpf.AgentConnTypeTKClose,
		throw:            true})[0]
	AssertConnEvent(t, connectEvent, ConnEventAssertions{
		expectPid:                  uint32(os.Getpid()),
		expectRemotePort:           80,
		expectLocalAddrFamily:      common.AF_INET,
		expectRemoteFamily:         common.AF_INET,
		expectReadBytesPredicator:  func(u uint64) bool { return u == 0 },
		expectWriteBytesPredicator: func(u uint64) bool { return u == 0 },
		expectLocalPort:            -1,
		expectConnEventType:        bpf.AgentConnTypeTKClose,
	})
}

func TestAccept(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	agentStopper := make(chan os.Signal, 1)

	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallAcceptEntry,
			bpf.AttachSyscallAcceptExit,
		},
		&connEventList,
		nil,
		nil,
		nil, agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()
	// ip, _ := common.GetIPAddrByInterfaceName("eth0")
	ip := "127.0.0.1"
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:      "hello\n",
		readResponse: true,
		writeSyscall: Write,
		readSyscall:  Read,
	})
	time.Sleep(1 * time.Second)
	if len(connEventList) == 0 {
		t.Fatalf("no conn event!")
	}

	connectEvent := findInterestedConnEvent(t, connEventList, FindInterestedConnEventOptions{
		findByLocalPort: true,
		findByConnType:  true,
		localPort:       uint16(echoTcpServerPort),
		connType:        bpf.AgentConnTypeTKConnect,
		throw:           true})[0]
	AssertConnEvent(t, connectEvent, ConnEventAssertions{
		expectPid:                  uint32(os.Getpid()),
		expectRemotePort:           -1,
		expectLocalAddrFamily:      common.AF_INET,
		expectRemoteFamily:         common.AF_INET,
		expectReadBytesPredicator:  func(u uint64) bool { return u == 0 },
		expectWriteBytesPredicator: func(u uint64) bool { return u == 0 },
		expectLocalPort:            echoTcpServerPort,
		expectConnEventType:        bpf.AgentConnTypeTKConnect,
	})
}

func TestRead(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	agentStopper := make(chan os.Signal)
	var connManager *conn.ConnManager = conn.InitConnManager()
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallReadEntry,
			bpf.AttachSyscallReadExit,
			bpf.AttachKProbeSecuritySocketRecvmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)

	defer func() {
		agentStopper <- MySignal{}
	}()
	ip := "127.0.0.1"
	sendMsg := "GET TestRead\n"
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:      sendMsg,
		readResponse: true,
		writeSyscall: Write,
		readSyscall:  Read,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	syscallEvent := syscallEvents[0]
	conn := connManager.FindConnection4(syscallEvent.SyscallEvent.Ke.ConnIdS.TgidFd)
	AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
		KernDataEventAssertConditions: KernDataEventAssertConditions{
			connIdDirect:     bpf.AgentTrafficDirectionTKIngress,
			pid:              uint64(os.Getpid()),
			fd:               uint32(conn.TgidFd),
			funcName:         "syscall",
			dataLen:          uint32(len(sendMsg)),
			seq:              1,
			step:             bpf.AgentStepTSYSCALL_IN,
			tsAssertFunction: func(u uint64) bool { return u > 0 },
		},
		bufSizeAssertFunction: func(u uint32) bool { return u == uint32(len(sendMsg)) },
		bufAssertFunction:     func(b []byte) bool { return string(b) == sendMsg },
	})
}

func TestRecvFrom(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	var connManager *conn.ConnManager = conn.InitConnManager()
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallRecvfromEntry,
			bpf.AttachSyscallRecvfromExit,
			bpf.AttachKProbeSecuritySocketRecvmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)

	defer func() {
		agentStopper <- MySignal{}
	}()
	ip := "127.0.0.1"
	sendMsg := "GET TestRecvFrom\n"
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:      sendMsg,
		readResponse: true,
		writeSyscall: Write,
		readSyscall:  RecvFrom,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	syscallEvent := syscallEvents[0]
	conn := connManager.FindConnection4(syscallEvent.SyscallEvent.Ke.ConnIdS.TgidFd)
	AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
		KernDataEventAssertConditions: KernDataEventAssertConditions{
			connIdDirect:     bpf.AgentTrafficDirectionTKIngress,
			pid:              uint64(os.Getpid()),
			fd:               uint32(conn.TgidFd),
			funcName:         "syscall",
			dataLen:          uint32(len(sendMsg)),
			seq:              1,
			step:             bpf.AgentStepTSYSCALL_IN,
			tsAssertFunction: func(u uint64) bool { return u > 0 },
		},
		bufSizeAssertFunction: func(u uint32) bool { return u == uint32(len(sendMsg)) },
		bufAssertFunction:     func(b []byte) bool { return string(b) == sendMsg },
	})
}

func TestReadv(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	agentStopper := make(chan os.Signal, 1)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	var connManager *conn.ConnManager = conn.InitConnManager()
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallReadvEntry,
			bpf.AttachSyscallReadvExit,
			bpf.AttachKProbeSecuritySocketRecvmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)

	defer func() {
		agentStopper <- MySignal{}
	}()
	ip := "127.0.0.1"
	sendMsg := "GET TestReadv\n"
	readBufSizeSlice := []int{len(sendMsg) / 2, len(sendMsg) - len(sendMsg)/2}
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:                t,
		server:           ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:          sendMsg,
		readResponse:     true,
		writeSyscall:     Write,
		readSyscall:      Readv,
		readBufSizeSlice: readBufSizeSlice,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	assert.Equal(t, len(readBufSizeSlice), len(syscallEvents))
	conn := connManager.FindConnection4(syscallEvents[0].SyscallEvent.Ke.ConnIdS.TgidFd)
	seq := uint64(1)
	for index, syscallEvent := range syscallEvents {
		AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
			KernDataEventAssertConditions: KernDataEventAssertConditions{
				connIdDirect:     bpf.AgentTrafficDirectionTKIngress,
				pid:              uint64(os.Getpid()),
				fd:               uint32(conn.TgidFd),
				funcName:         "syscall",
				dataLen:          uint32(readBufSizeSlice[index]),
				seq:              seq,
				step:             bpf.AgentStepTSYSCALL_IN,
				tsAssertFunction: func(u uint64) bool { return u > 0 },
			},
			bufSizeAssertFunction: func(u uint32) bool { return u == uint32(readBufSizeSlice[index]) },
			bufAssertFunction:     func(b []byte) bool { return string(b) == sendMsg[seq-1:int(seq)-1+readBufSizeSlice[index]] },
		})
		seq += uint64(readBufSizeSlice[index])
	}
}

func TestRecvmsg(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	agentStopper := make(chan os.Signal, 1)
	var connManager *conn.ConnManager = conn.InitConnManager()
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallRecvMsgEntry,
			bpf.AttachSyscallRecvMsgExit,
			bpf.AttachKProbeSecuritySocketRecvmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()

	ip := "127.0.0.1"
	sendMsg := "GET TestRecvmsg\n"
	readBufSizeSlice := []int{len(sendMsg) / 2, len(sendMsg) - len(sendMsg)/2}
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:                t,
		server:           ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:          sendMsg,
		readResponse:     true,
		writeSyscall:     Write,
		readSyscall:      Recvmsg,
		readBufSizeSlice: readBufSizeSlice,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	assert.Equal(t, len(readBufSizeSlice), len(syscallEvents))
	conn := connManager.FindConnection4(syscallEvents[0].SyscallEvent.Ke.ConnIdS.TgidFd)
	seq := uint64(1)
	for index, syscallEvent := range syscallEvents {
		AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
			KernDataEventAssertConditions: KernDataEventAssertConditions{
				connIdDirect:     bpf.AgentTrafficDirectionTKIngress,
				pid:              uint64(os.Getpid()),
				fd:               uint32(conn.TgidFd),
				funcName:         "syscall",
				dataLen:          uint32(readBufSizeSlice[index]),
				seq:              seq,
				step:             bpf.AgentStepTSYSCALL_IN,
				tsAssertFunction: func(u uint64) bool { return u > 0 },
			},
			bufSizeAssertFunction: func(u uint32) bool { return u == uint32(readBufSizeSlice[index]) },
			bufAssertFunction:     func(b []byte) bool { return string(b) == sendMsg[seq-1:int(seq)-1+readBufSizeSlice[index]] },
		})
		seq += uint64(readBufSizeSlice[index])
	}
}

func TestWrite(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	var connManager *conn.ConnManager = conn.InitConnManager()
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallWriteEntry,
			bpf.AttachSyscallWriteExit,
			bpf.AttachKProbeSecuritySocketSendmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)

	defer func() {
		agentStopper <- MySignal{}
	}()
	ip := "127.0.0.1"
	sendMsg := "GET TestWrite\n"
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:      sendMsg,
		readResponse: true,
		writeSyscall: Write,
		readSyscall:  Read,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	syscallEvent := syscallEvents[0]
	conn := connManager.FindConnection4(syscallEvent.SyscallEvent.Ke.ConnIdS.TgidFd)
	AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
		KernDataEventAssertConditions: KernDataEventAssertConditions{
			connIdDirect:     bpf.AgentTrafficDirectionTKEgress,
			pid:              uint64(os.Getpid()),
			fd:               uint32(conn.TgidFd),
			funcName:         "syscall",
			dataLen:          uint32(len(sendMsg)),
			seq:              1,
			step:             bpf.AgentStepTSYSCALL_OUT,
			tsAssertFunction: func(u uint64) bool { return u > 0 },
		},
		bufSizeAssertFunction: func(u uint32) bool { return u == uint32(len(sendMsg)) },
		bufAssertFunction:     func(b []byte) bool { return string(b) == sendMsg },
	})
}

func TestSendto(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	var connManager *conn.ConnManager = conn.InitConnManager()
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallSendtoEntry,
			bpf.AttachSyscallSendtoExit,
			bpf.AttachKProbeSecuritySocketSendmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()

	ip := "127.0.0.1"
	sendMsg := "GET TestSendto\n"
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:      sendMsg,
		readResponse: true,
		writeSyscall: SentTo,
		readSyscall:  Read,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	syscallEvent := syscallEvents[0]
	conn := connManager.FindConnection4(syscallEvent.SyscallEvent.Ke.ConnIdS.TgidFd)
	AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
		KernDataEventAssertConditions: KernDataEventAssertConditions{
			connIdDirect:     bpf.AgentTrafficDirectionTKEgress,
			pid:              uint64(os.Getpid()),
			fd:               uint32(conn.TgidFd),
			funcName:         "syscall",
			dataLen:          uint32(len(sendMsg)),
			seq:              1,
			step:             bpf.AgentStepTSYSCALL_OUT,
			tsAssertFunction: func(u uint64) bool { return u > 0 },
		},
		bufSizeAssertFunction: func(u uint32) bool { return u == uint32(len(sendMsg)) },
		bufAssertFunction:     func(b []byte) bool { return string(b) == sendMsg },
	})
}

func TestWritev(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	agentStopper := make(chan os.Signal, 1)
	var connManager *conn.ConnManager = conn.InitConnManager()
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallWritevEntry,
			bpf.AttachSyscallWritevExit,
			bpf.AttachKProbeSecuritySocketSendmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()

	ip := "127.0.0.1"
	sendMessages := []string{"GET writevhellohellohellohello\n", "abchellohellohellohello\n"}
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		messageSlice: sendMessages,
		readResponse: true,
		writeSyscall: Writev,
		readSyscall:  Read,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	assert.Equal(t, 2, len(syscallEvents))
	conn := connManager.FindConnection4(syscallEvents[0].SyscallEvent.Ke.ConnIdS.TgidFd)
	seq := uint64(1)
	for index, syscallEvent := range syscallEvents {
		AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
			KernDataEventAssertConditions: KernDataEventAssertConditions{
				connIdDirect:     bpf.AgentTrafficDirectionTKEgress,
				pid:              uint64(os.Getpid()),
				fd:               uint32(conn.TgidFd),
				funcName:         "syscall",
				dataLen:          uint32(len(sendMessages[index])),
				seq:              seq,
				step:             bpf.AgentStepTSYSCALL_OUT,
				tsAssertFunction: func(u uint64) bool { return u > 0 },
			},
			bufSizeAssertFunction: func(u uint32) bool { return u == uint32(len(sendMessages[index])) },
			bufAssertFunction:     func(b []byte) bool { return string(b) == sendMessages[index] },
		})
		seq += uint64(len(sendMessages[index]))
	}
}

func TestSendMsg(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	var connManager *conn.ConnManager = conn.InitConnManager()
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallSendMsgEntry,
			bpf.AttachSyscallSendMsgExit,
			bpf.AttachKProbeSecuritySocketSendmsgEntry,
		},
		&connEventList,
		&syscallEventList,
		nil,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)
	defer func() {
		agentStopper <- MySignal{}
	}()

	ip := "127.0.0.1"
	sendMessages := []string{"GET sendmsghellohellohellohello\n", "aabchellohellohellohello\n"}
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		messageSlice: sendMessages,
		readResponse: true,
		writeSyscall: Sendmsg,
		readSyscall:  Read,
	})
	time.Sleep(500 * time.Millisecond)

	assert.NotEmpty(t, syscallEventList)
	syscallEvents := findInterestedSyscallEvents(t, syscallEventList, FindInterestedSyscallEventOptions{
		findByRemotePort: true,
		remotePort:       uint16(echoTcpServerPort),
		connEventList:    connEventList,
	})
	assert.Equal(t, 2, len(syscallEvents))
	conn := connManager.FindConnection4(syscallEvents[0].SyscallEvent.Ke.ConnIdS.TgidFd)
	seq := uint64(1)
	for index, syscallEvent := range syscallEvents {
		AssertSyscallEventData(t, syscallEvent, SyscallDataEventAssertConditions{
			KernDataEventAssertConditions: KernDataEventAssertConditions{connIdDirect: bpf.AgentTrafficDirectionTKEgress,
				pid:              uint64(os.Getpid()),
				fd:               uint32(conn.TgidFd),
				funcName:         "syscall",
				dataLen:          uint32(len(sendMessages[index])),
				seq:              seq,
				step:             bpf.AgentStepTSYSCALL_OUT,
				tsAssertFunction: func(u uint64) bool { return u > 0 },
			},
			bufSizeAssertFunction: func(u uint32) bool { return u == uint32(len(sendMessages[index])) },
			bufAssertFunction:     func(b []byte) bool { return string(b) == sendMessages[index] },
		})
		seq += uint64(len(sendMessages[index]))
	}
}

func TestIpXmit(t *testing.T) {
	StartEchoTcpServerAndWait()
	connEventList := make([]bpf.AgentConnEvtT, 0)
	syscallEventList := make([]bpf.SyscallEventData, 0)
	kernEventList := make([]bpf.AgentKernEvt, 0)
	var connManager *conn.ConnManager = conn.InitConnManager()
	agentStopper := make(chan os.Signal, 1)
	StartAgent(
		[]bpf.AttachBpfProgFunction{
			bpf.AttachSyscallConnectEntry,
			bpf.AttachSyscallConnectExit,
			bpf.AttachSyscallWriteEntry,
			bpf.AttachSyscallWriteExit,
			bpf.AttachKProbeSecuritySocketSendmsgEntry,
			bpf.AttachKProbeIpQueueXmitEntry,
		},
		&connEventList,
		&syscallEventList,
		&kernEventList,
		func(cm *conn.ConnManager) {
			*connManager = *cm
		}, agentStopper)

	defer func() {
		agentStopper <- MySignal{}
	}()
	ip := "127.0.0.1"
	sendMsg := "GET TestIpXmit\n"
	WriteToEchoTcpServerAndReadResponse(WriteToEchoServerOptions{
		t:            t,
		server:       ip + ":" + fmt.Sprint(echoTcpServerPort),
		message:      sendMsg,
		readResponse: true,
		writeSyscall: Write,
		readSyscall:  Read,
	})
	time.Sleep(500 * time.Millisecond)
	intersetedKernEvents := findInterestedKernEvents(t, kernEventList, FindInterestedKernEventOptions{
		connEventList:          connEventList,
		findDataLenGtZeroEvent: true,
		findByDirect:           true,
		direct:                 bpf.AgentTrafficDirectionTKEgress,
	})
	assert.Equal(t, 1, len(intersetedKernEvents))
	kernEvent := intersetedKernEvents[0]
	conn := connManager.FindConnection4(kernEvent.ConnIdS.TgidFd)
	AssertKernEvent(t, &kernEvent, KernDataEventAssertConditions{

		connIdDirect:     bpf.AgentTrafficDirectionTKEgress,
		pid:              uint64(os.Getpid()),
		fd:               uint32(conn.TgidFd),
		funcName:         "ip_queue_xmit",
		dataLen:          uint32(len(sendMsg)),
		seq:              1,
		step:             bpf.AgentStepTIP_OUT,
		tsAssertFunction: func(u uint64) bool { return u > 0 },
	})
}