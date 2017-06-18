package liquid

import (
	"reflect"

	"github.com/ybbus/jsonrpc"
)

// RPCClientEngine connects via JSON-RPC to a Liquid template server.
type RPCClientEngine struct {
	rpcClient    *jsonrpc.RPCClient
	rpcSessionID string
}

// DefaultServer is the default HTTP address for a Liquid template server.
// This is an unclaimed port number from https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Registered_ports
const DefaultServer = "localhost:4545"

type remoteTemplate struct {
	engine RemoteEngine
	text   []byte
}

// NewRPCClientEngine creates a RemoteEngine.
func NewRPCClientEngine(address string) RemoteEngine {
	rpcClient := jsonrpc.NewRPCClient("http://" + address)
	return &RPCClientEngine{rpcClient: rpcClient}
}

// Parse parses the template.
func (engine *RPCClientEngine) Parse(text []byte) (Template, error) {
	return &remoteTemplate{engine, text}, nil
}

// RPCError wraps jsonrpc.RPCError into an Error
type RPCError struct{ jsonrpc.RPCError }

func (engine *RPCError) Error() string {
	return engine.Message
}

func (engine *RPCClientEngine) getSessionID() string {
	if engine.rpcSessionID != "" {
		return engine.rpcSessionID
	}
	res, err := engine.rpcClient.Call("session")
	if err != nil {
		panic(err)
	}
	if res.Error != nil {
		panic(&RPCError{*res.Error})
	}
	var result struct {
		SessionID string
	}
	res.GetObject(&result)
	engine.rpcSessionID = result.SessionID
	return engine.rpcSessionID
}

func (engine *RPCClientEngine) rpcCall(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	args := append([]interface{}{engine.getSessionID()}, params...)
	res, err := engine.rpcClient.Call(method, args...)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, &RPCError{*res.Error}
	}
	return res, nil
}

func (engine *RPCClientEngine) FileUrlMap(m map[string]string) {
	_, err := engine.rpcCall("fileUrls", m)
	if err != nil {
		panic(err)
	}
}

func (engine *RPCClientEngine) IncludeDirs(dirs []string) {
	_, err := engine.rpcCall("includeDirs", dirs)
	if err != nil {
		panic(err)
	}
}

// ParseAndRender parses and then renders the template.
func (engine *RPCClientEngine) ParseAndRender(text []byte, scope map[string]interface{}) ([]byte, error) {
	scope = prepForJSON(scope).(map[string]interface{})

	res, err := engine.rpcCall("render", string(text), scope)
	if err != nil {
		return nil, err
	}

	var render struct {
		Text string
	}
	err = res.GetObject(&render)
	if err != nil {
		return nil, err
	}
	return []byte(render.Text), nil
}

// Render renders the template.
func (template *remoteTemplate) Render(scope map[string]interface{}) ([]byte, error) {
	return template.engine.ParseAndRender(template.text, scope)
}

func prepForJSON(value interface{}) interface{} {
	ref := reflect.ValueOf(value)
	switch ref.Kind() {
	case reflect.Map:
		m := map[string]interface{}{}
		for _, k := range ref.MapKeys() {
			m[k.String()] = prepForJSON(ref.MapIndex(k).Interface())
		}
		return m
	case reflect.Slice:
		s := make([]interface{}, ref.Len())
		for i := 0; i < ref.Len(); i++ {
			s[i] = prepForJSON(ref.Index(i).Interface())
		}
		return s
	default:
		return value
	}
}
