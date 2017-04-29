package knx

import (
	"bytes"
	"errors"
)

// IPv4 address
type Address [4]byte

// Port number
type Port uint16

// Host information
type HostInfo struct {
	Address Address
	Port    Port
}

var hostInfoPreamble = [2]byte{8, 1}

func (info HostInfo) writeTo(w *bytes.Buffer) error {
	return writeSequence(w, hostInfoPreamble, info)
}

func readHostInfo(r *bytes.Reader) (*HostInfo, error) {
	var length, proto byte
	info := &HostInfo{}

	err := readSequence(r, &length, &proto, info)
	if err != nil { return nil, err }

	if length != 8 {
		return nil, errors.New("Host info structure length is invalid")
	}

	if proto != 1 {
		return nil, errors.New("Host info protocol is unknown")
	}

	return info, nil
}

// Connection request
type ConnectionRequest struct {
	Control HostInfo
	Tunnel  HostInfo
}

func (req ConnectionRequest) describe() (serviceIdent, int) {
	return connectionRequestService, 20
}

var connReqInfo = [4]byte{4, 4, 2, 0}

func (req ConnectionRequest) writeTo(w *bytes.Buffer) error {
	err := req.Control.writeTo(w)
	if err != nil { return err }

	err = req.Tunnel.writeTo(w)
	if err != nil { return err }

	_, err = w.Write(connReqInfo[:])
	return err
}

// Connection response
type ConnectionResponse struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func readConnectionResponse(r *bytes.Reader) (*ConnectionResponse, error) {
	var channel, status byte

	err := readSequence(r, &channel, &status)
	if err != nil { return nil, err }

	host, err := readHostInfo(r)
	if err != nil { return nil, err }

	return &ConnectionResponse{channel, status, *host}, nil
}

// Connection state request
type ConnectionStateRequest struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func (req ConnectionStateRequest) describe() (serviceIdent, int) {
	return connectionStateRequestService, 10
}

func (req ConnectionStateRequest) writeTo(w *bytes.Buffer) error {
	err := writeSequence(w, req.Channel, req.Status)
	if err != nil { return err }

	return req.Host.writeTo(w)
}

// Connection state response
type ConnectionStateResponse struct {
	Channel byte
	Status  byte
}

func readConnectionStateResponse(r *bytes.Reader) (*ConnectionStateResponse, error) {
	res := &ConnectionStateResponse{}
	return res, readSequence(r, &res.Channel, &res.Status)
}

// Disconnect request
type DisconnectRequest struct {
	Channel byte
	Status  byte
	Host    HostInfo
}

func readDisconnectRequest(r *bytes.Reader) (*DisconnectRequest, error) {
	var channel, status byte

	err := readSequence(r, &channel, &status)
	if err != nil { return nil, err }

	host, err := readHostInfo(r)
	if err != nil { return nil, err }

	return &DisconnectRequest{channel, status, *host}, nil
}

func (req DisconnectRequest) describe() (serviceIdent, int) {
	return disconnectRequestService, 10
}

func (req DisconnectRequest) writeTo(w *bytes.Buffer) error {
	err := writeSequence(w, req.Channel, req.Status)
	if err != nil { return err }

	return req.Host.writeTo(w)
}