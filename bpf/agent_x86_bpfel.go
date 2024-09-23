// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64

package bpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type AgentConnEvtT struct {
	ConnInfo AgentConnInfoT
	ConnType AgentConnTypeT
	_        [4]byte
	Ts       uint64
}

type AgentConnIdS_t struct {
	TgidFd  uint64
	NoTrace bool
	_       [7]byte
}

type AgentConnInfoT struct {
	ConnId struct {
		Upid struct {
			Pid            uint32
			_              [4]byte
			StartTimeTicks uint64
		}
		Fd   int32
		_    [4]byte
		Tsid uint64
	}
	ReadBytes  uint64
	WriteBytes uint64
	Laddr      struct {
		In6 struct {
			Sin6Family   uint16
			Sin6Port     uint16
			Sin6Flowinfo uint32
			Sin6Addr     struct{ In6U struct{ U6Addr8 [16]uint8 } }
			Sin6ScopeId  uint32
		}
	}
	Raddr struct {
		In6 struct {
			Sin6Family   uint16
			Sin6Port     uint16
			Sin6Flowinfo uint32
			Sin6Addr     struct{ In6U struct{ U6Addr8 [16]uint8 } }
			Sin6ScopeId  uint32
		}
	}
	Protocol            AgentTrafficProtocolT
	Role                AgentEndpointRoleT
	PrevCount           uint64
	PrevBuf             [4]int8
	PrependLengthHeader bool
	NoTrace             bool
	Ssl                 bool
	_                   [1]byte
}

type AgentConnTypeT uint32

const (
	AgentConnTypeTKConnect       AgentConnTypeT = 0
	AgentConnTypeTKClose         AgentConnTypeT = 1
	AgentConnTypeTKProtocolInfer AgentConnTypeT = 2
)

type AgentControlValueIndexT uint32

const (
	AgentControlValueIndexTKTargetTGIDIndex   AgentControlValueIndexT = 0
	AgentControlValueIndexTKStirlingTGIDIndex AgentControlValueIndexT = 1
	AgentControlValueIndexTKEnabledXdpIndex   AgentControlValueIndexT = 2
	AgentControlValueIndexTKNumControlValues  AgentControlValueIndexT = 3
)

type AgentEndpointRoleT uint32

const (
	AgentEndpointRoleTKRoleClient  AgentEndpointRoleT = 1
	AgentEndpointRoleTKRoleServer  AgentEndpointRoleT = 2
	AgentEndpointRoleTKRoleUnknown AgentEndpointRoleT = 4
)

type AgentKernEvt struct {
	FuncName [16]int8
	Ts       uint64
	Seq      uint64
	Len      uint32
	Flags    uint8
	_        [3]byte
	ConnIdS  AgentConnIdS_t
	IsSample int32
	Step     AgentStepT
}

type AgentKernEvtData struct {
	Ke      AgentKernEvt
	BufSize uint32
	Msg     [30720]int8
	_       [4]byte
}

type AgentSockKey struct {
	Sip   [2]uint64
	Dip   [2]uint64
	Sport uint16
	Dport uint16
	_     [4]byte
}

type AgentStepT uint32

const (
	AgentStepTStart       AgentStepT = 0
	AgentStepTSYSCALL_OUT AgentStepT = 1
	AgentStepTTCP_OUT     AgentStepT = 2
	AgentStepTIP_OUT      AgentStepT = 3
	AgentStepTQDISC_OUT   AgentStepT = 4
	AgentStepTDEV_OUT     AgentStepT = 5
	AgentStepTNIC_OUT     AgentStepT = 6
	AgentStepTNIC_IN      AgentStepT = 7
	AgentStepTDEV_IN      AgentStepT = 8
	AgentStepTIP_IN       AgentStepT = 9
	AgentStepTTCP_IN      AgentStepT = 10
	AgentStepTUSER_COPY   AgentStepT = 11
	AgentStepTSYSCALL_IN  AgentStepT = 12
	AgentStepTEnd         AgentStepT = 13
)

type AgentTrafficDirectionT uint32

const (
	AgentTrafficDirectionTKEgress  AgentTrafficDirectionT = 0
	AgentTrafficDirectionTKIngress AgentTrafficDirectionT = 1
)

type AgentTrafficProtocolT uint32

const (
	AgentTrafficProtocolTKProtocolUnset   AgentTrafficProtocolT = 0
	AgentTrafficProtocolTKProtocolUnknown AgentTrafficProtocolT = 1
	AgentTrafficProtocolTKProtocolHTTP    AgentTrafficProtocolT = 2
	AgentTrafficProtocolTKProtocolHTTP2   AgentTrafficProtocolT = 3
	AgentTrafficProtocolTKProtocolMySQL   AgentTrafficProtocolT = 4
	AgentTrafficProtocolTKProtocolCQL     AgentTrafficProtocolT = 5
	AgentTrafficProtocolTKProtocolPGSQL   AgentTrafficProtocolT = 6
	AgentTrafficProtocolTKProtocolDNS     AgentTrafficProtocolT = 7
	AgentTrafficProtocolTKProtocolRedis   AgentTrafficProtocolT = 8
	AgentTrafficProtocolTKProtocolNATS    AgentTrafficProtocolT = 9
	AgentTrafficProtocolTKProtocolMongo   AgentTrafficProtocolT = 10
	AgentTrafficProtocolTKProtocolKafka   AgentTrafficProtocolT = 11
	AgentTrafficProtocolTKProtocolMux     AgentTrafficProtocolT = 12
	AgentTrafficProtocolTKProtocolAMQP    AgentTrafficProtocolT = 13
	AgentTrafficProtocolTKNumProtocols    AgentTrafficProtocolT = 14
)

// LoadAgent returns the embedded CollectionSpec for Agent.
func LoadAgent() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_AgentBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load Agent: %w", err)
	}

	return spec, err
}

// LoadAgentObjects loads Agent and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*AgentObjects
//	*AgentPrograms
//	*AgentMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func LoadAgentObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := LoadAgent()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// AgentSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type AgentSpecs struct {
	AgentProgramSpecs
	AgentMapSpecs
}

// AgentSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type AgentProgramSpecs struct {
	DevHardStartXmit                   *ebpf.ProgramSpec `ebpf:"dev_hard_start_xmit"`
	DevQueueXmit                       *ebpf.ProgramSpec `ebpf:"dev_queue_xmit"`
	IpQueueXmit                        *ebpf.ProgramSpec `ebpf:"ip_queue_xmit"`
	IpQueueXmit2                       *ebpf.ProgramSpec `ebpf:"ip_queue_xmit2"`
	IpRcvCore                          *ebpf.ProgramSpec `ebpf:"ip_rcv_core"`
	SecuritySocketRecvmsgEnter         *ebpf.ProgramSpec `ebpf:"security_socket_recvmsg_enter"`
	SecuritySocketSendmsgEnter         *ebpf.ProgramSpec `ebpf:"security_socket_sendmsg_enter"`
	SkbCopyDatagramIovec               *ebpf.ProgramSpec `ebpf:"skb_copy_datagram_iovec"`
	SkbCopyDatagramIter                *ebpf.ProgramSpec `ebpf:"skb_copy_datagram_iter"`
	SockAllocRet                       *ebpf.ProgramSpec `ebpf:"sock_alloc_ret"`
	TcpDestroySock                     *ebpf.ProgramSpec `ebpf:"tcp_destroy_sock"`
	TcpQueueRcv                        *ebpf.ProgramSpec `ebpf:"tcp_queue_rcv"`
	TcpRcvEstablished                  *ebpf.ProgramSpec `ebpf:"tcp_rcv_established"`
	TcpV4DoRcv                         *ebpf.ProgramSpec `ebpf:"tcp_v4_do_rcv"`
	TcpV4Rcv                           *ebpf.ProgramSpec `ebpf:"tcp_v4_rcv"`
	TracepointNetifReceiveSkb          *ebpf.ProgramSpec `ebpf:"tracepoint__netif_receive_skb"`
	TracepointSyscallsSysEnterAccept4  *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_accept4"`
	TracepointSyscallsSysEnterClose    *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_close"`
	TracepointSyscallsSysEnterConnect  *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_connect"`
	TracepointSyscallsSysEnterRead     *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_read"`
	TracepointSyscallsSysEnterReadv    *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_readv"`
	TracepointSyscallsSysEnterRecvfrom *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_recvfrom"`
	TracepointSyscallsSysEnterRecvmsg  *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_recvmsg"`
	TracepointSyscallsSysEnterSendmsg  *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_sendmsg"`
	TracepointSyscallsSysEnterSendto   *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_sendto"`
	TracepointSyscallsSysEnterWrite    *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_write"`
	TracepointSyscallsSysEnterWritev   *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_enter_writev"`
	TracepointSyscallsSysExitAccept4   *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_accept4"`
	TracepointSyscallsSysExitClose     *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_close"`
	TracepointSyscallsSysExitConnect   *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_connect"`
	TracepointSyscallsSysExitRead      *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_read"`
	TracepointSyscallsSysExitReadv     *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_readv"`
	TracepointSyscallsSysExitRecvfrom  *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_recvfrom"`
	TracepointSyscallsSysExitRecvmsg   *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_recvmsg"`
	TracepointSyscallsSysExitSendmsg   *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_sendmsg"`
	TracepointSyscallsSysExitSendto    *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_sendto"`
	TracepointSyscallsSysExitWrite     *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_write"`
	TracepointSyscallsSysExitWritev    *ebpf.ProgramSpec `ebpf:"tracepoint__syscalls__sys_exit_writev"`
	XdpProxy                           *ebpf.ProgramSpec `ebpf:"xdp_proxy"`
}

// AgentMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type AgentMapSpecs struct {
	AcceptArgsMap        *ebpf.MapSpec `ebpf:"accept_args_map"`
	ActiveSslReadArgsMap *ebpf.MapSpec `ebpf:"active_ssl_read_args_map"`
	CloseArgsMap         *ebpf.MapSpec `ebpf:"close_args_map"`
	ConnEvtRb            *ebpf.MapSpec `ebpf:"conn_evt_rb"`
	ConnInfoMap          *ebpf.MapSpec `ebpf:"conn_info_map"`
	ConnInfoT_map        *ebpf.MapSpec `ebpf:"conn_info_t_map"`
	ConnectArgsMap       *ebpf.MapSpec `ebpf:"connect_args_map"`
	ControlValues        *ebpf.MapSpec `ebpf:"control_values"`
	EnabledLocalIpv4Map  *ebpf.MapSpec `ebpf:"enabled_local_ipv4_map"`
	EnabledLocalPortMap  *ebpf.MapSpec `ebpf:"enabled_local_port_map"`
	EnabledRemoteIpv4Map *ebpf.MapSpec `ebpf:"enabled_remote_ipv4_map"`
	EnabledRemotePortMap *ebpf.MapSpec `ebpf:"enabled_remote_port_map"`
	Rb                   *ebpf.MapSpec `ebpf:"rb"`
	ReadArgsMap          *ebpf.MapSpec `ebpf:"read_args_map"`
	SockKeyConnIdMap     *ebpf.MapSpec `ebpf:"sock_key_conn_id_map"`
	SockXmitMap          *ebpf.MapSpec `ebpf:"sock_xmit_map"`
	SslUserSpaceCallMap  *ebpf.MapSpec `ebpf:"ssl_user_space_call_map"`
	SyscallDataMap       *ebpf.MapSpec `ebpf:"syscall_data_map"`
	SyscallRb            *ebpf.MapSpec `ebpf:"syscall_rb"`
	WriteArgsMap         *ebpf.MapSpec `ebpf:"write_args_map"`
}

// AgentObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to LoadAgentObjects or ebpf.CollectionSpec.LoadAndAssign.
type AgentObjects struct {
	AgentPrograms
	AgentMaps
}

func (o *AgentObjects) Close() error {
	return _AgentClose(
		&o.AgentPrograms,
		&o.AgentMaps,
	)
}

// AgentMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to LoadAgentObjects or ebpf.CollectionSpec.LoadAndAssign.
type AgentMaps struct {
	AcceptArgsMap        *ebpf.Map `ebpf:"accept_args_map"`
	ActiveSslReadArgsMap *ebpf.Map `ebpf:"active_ssl_read_args_map"`
	CloseArgsMap         *ebpf.Map `ebpf:"close_args_map"`
	ConnEvtRb            *ebpf.Map `ebpf:"conn_evt_rb"`
	ConnInfoMap          *ebpf.Map `ebpf:"conn_info_map"`
	ConnInfoT_map        *ebpf.Map `ebpf:"conn_info_t_map"`
	ConnectArgsMap       *ebpf.Map `ebpf:"connect_args_map"`
	ControlValues        *ebpf.Map `ebpf:"control_values"`
	EnabledLocalIpv4Map  *ebpf.Map `ebpf:"enabled_local_ipv4_map"`
	EnabledLocalPortMap  *ebpf.Map `ebpf:"enabled_local_port_map"`
	EnabledRemoteIpv4Map *ebpf.Map `ebpf:"enabled_remote_ipv4_map"`
	EnabledRemotePortMap *ebpf.Map `ebpf:"enabled_remote_port_map"`
	Rb                   *ebpf.Map `ebpf:"rb"`
	ReadArgsMap          *ebpf.Map `ebpf:"read_args_map"`
	SockKeyConnIdMap     *ebpf.Map `ebpf:"sock_key_conn_id_map"`
	SockXmitMap          *ebpf.Map `ebpf:"sock_xmit_map"`
	SslUserSpaceCallMap  *ebpf.Map `ebpf:"ssl_user_space_call_map"`
	SyscallDataMap       *ebpf.Map `ebpf:"syscall_data_map"`
	SyscallRb            *ebpf.Map `ebpf:"syscall_rb"`
	WriteArgsMap         *ebpf.Map `ebpf:"write_args_map"`
}

func (m *AgentMaps) Close() error {
	return _AgentClose(
		m.AcceptArgsMap,
		m.ActiveSslReadArgsMap,
		m.CloseArgsMap,
		m.ConnEvtRb,
		m.ConnInfoMap,
		m.ConnInfoT_map,
		m.ConnectArgsMap,
		m.ControlValues,
		m.EnabledLocalIpv4Map,
		m.EnabledLocalPortMap,
		m.EnabledRemoteIpv4Map,
		m.EnabledRemotePortMap,
		m.Rb,
		m.ReadArgsMap,
		m.SockKeyConnIdMap,
		m.SockXmitMap,
		m.SslUserSpaceCallMap,
		m.SyscallDataMap,
		m.SyscallRb,
		m.WriteArgsMap,
	)
}

// AgentPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to LoadAgentObjects or ebpf.CollectionSpec.LoadAndAssign.
type AgentPrograms struct {
	DevHardStartXmit                   *ebpf.Program `ebpf:"dev_hard_start_xmit"`
	DevQueueXmit                       *ebpf.Program `ebpf:"dev_queue_xmit"`
	IpQueueXmit                        *ebpf.Program `ebpf:"ip_queue_xmit"`
	IpQueueXmit2                       *ebpf.Program `ebpf:"ip_queue_xmit2"`
	IpRcvCore                          *ebpf.Program `ebpf:"ip_rcv_core"`
	SecuritySocketRecvmsgEnter         *ebpf.Program `ebpf:"security_socket_recvmsg_enter"`
	SecuritySocketSendmsgEnter         *ebpf.Program `ebpf:"security_socket_sendmsg_enter"`
	SkbCopyDatagramIovec               *ebpf.Program `ebpf:"skb_copy_datagram_iovec"`
	SkbCopyDatagramIter                *ebpf.Program `ebpf:"skb_copy_datagram_iter"`
	SockAllocRet                       *ebpf.Program `ebpf:"sock_alloc_ret"`
	TcpDestroySock                     *ebpf.Program `ebpf:"tcp_destroy_sock"`
	TcpQueueRcv                        *ebpf.Program `ebpf:"tcp_queue_rcv"`
	TcpRcvEstablished                  *ebpf.Program `ebpf:"tcp_rcv_established"`
	TcpV4DoRcv                         *ebpf.Program `ebpf:"tcp_v4_do_rcv"`
	TcpV4Rcv                           *ebpf.Program `ebpf:"tcp_v4_rcv"`
	TracepointNetifReceiveSkb          *ebpf.Program `ebpf:"tracepoint__netif_receive_skb"`
	TracepointSyscallsSysEnterAccept4  *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_accept4"`
	TracepointSyscallsSysEnterClose    *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_close"`
	TracepointSyscallsSysEnterConnect  *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_connect"`
	TracepointSyscallsSysEnterRead     *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_read"`
	TracepointSyscallsSysEnterReadv    *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_readv"`
	TracepointSyscallsSysEnterRecvfrom *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_recvfrom"`
	TracepointSyscallsSysEnterRecvmsg  *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_recvmsg"`
	TracepointSyscallsSysEnterSendmsg  *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_sendmsg"`
	TracepointSyscallsSysEnterSendto   *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_sendto"`
	TracepointSyscallsSysEnterWrite    *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_write"`
	TracepointSyscallsSysEnterWritev   *ebpf.Program `ebpf:"tracepoint__syscalls__sys_enter_writev"`
	TracepointSyscallsSysExitAccept4   *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_accept4"`
	TracepointSyscallsSysExitClose     *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_close"`
	TracepointSyscallsSysExitConnect   *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_connect"`
	TracepointSyscallsSysExitRead      *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_read"`
	TracepointSyscallsSysExitReadv     *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_readv"`
	TracepointSyscallsSysExitRecvfrom  *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_recvfrom"`
	TracepointSyscallsSysExitRecvmsg   *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_recvmsg"`
	TracepointSyscallsSysExitSendmsg   *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_sendmsg"`
	TracepointSyscallsSysExitSendto    *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_sendto"`
	TracepointSyscallsSysExitWrite     *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_write"`
	TracepointSyscallsSysExitWritev    *ebpf.Program `ebpf:"tracepoint__syscalls__sys_exit_writev"`
	XdpProxy                           *ebpf.Program `ebpf:"xdp_proxy"`
}

func (p *AgentPrograms) Close() error {
	return _AgentClose(
		p.DevHardStartXmit,
		p.DevQueueXmit,
		p.IpQueueXmit,
		p.IpQueueXmit2,
		p.IpRcvCore,
		p.SecuritySocketRecvmsgEnter,
		p.SecuritySocketSendmsgEnter,
		p.SkbCopyDatagramIovec,
		p.SkbCopyDatagramIter,
		p.SockAllocRet,
		p.TcpDestroySock,
		p.TcpQueueRcv,
		p.TcpRcvEstablished,
		p.TcpV4DoRcv,
		p.TcpV4Rcv,
		p.TracepointNetifReceiveSkb,
		p.TracepointSyscallsSysEnterAccept4,
		p.TracepointSyscallsSysEnterClose,
		p.TracepointSyscallsSysEnterConnect,
		p.TracepointSyscallsSysEnterRead,
		p.TracepointSyscallsSysEnterReadv,
		p.TracepointSyscallsSysEnterRecvfrom,
		p.TracepointSyscallsSysEnterRecvmsg,
		p.TracepointSyscallsSysEnterSendmsg,
		p.TracepointSyscallsSysEnterSendto,
		p.TracepointSyscallsSysEnterWrite,
		p.TracepointSyscallsSysEnterWritev,
		p.TracepointSyscallsSysExitAccept4,
		p.TracepointSyscallsSysExitClose,
		p.TracepointSyscallsSysExitConnect,
		p.TracepointSyscallsSysExitRead,
		p.TracepointSyscallsSysExitReadv,
		p.TracepointSyscallsSysExitRecvfrom,
		p.TracepointSyscallsSysExitRecvmsg,
		p.TracepointSyscallsSysExitSendmsg,
		p.TracepointSyscallsSysExitSendto,
		p.TracepointSyscallsSysExitWrite,
		p.TracepointSyscallsSysExitWritev,
		p.XdpProxy,
	)
}

func _AgentClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed agent_x86_bpfel.o
var _AgentBytes []byte
