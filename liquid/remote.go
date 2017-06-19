package liquid

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/ybbus/jsonrpc"
)

// RPCClientEngine connects via JSON-RPC to a Liquid template server.
type RPCClientEngine struct {
	rpcClient *jsonrpc.RPCClient
}

// DefaultServer is the default HTTP address for a Liquid template server.
// This is an unclaimed port number from https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers#Registered_ports
const DefaultServer = "localhost:4545"

// RPCVersion is the Liquid Template Server RPC version.
const RPCVersion = "0.0.1"

type remoteTemplate struct {
	engine RemoteEngine
	text   []byte
}

// RPCError wraps jsonrpc.RPCError into an Error.
type RPCError struct{ jsonrpc.RPCError }

func (engine *RPCError) Error() string {
	return engine.Message
}

// NewRPCClientEngine creates a RemoteEngine.
func NewRPCClientEngine(address string) (RemoteEngine, error) {
	rpcClient := jsonrpc.NewRPCClient("http://" + address)
	engine := RPCClientEngine{rpcClient: rpcClient}
	if err := engine.createSession(); err != nil {
		return nil, err
	}
	return &engine, nil
}

// Parse parses the template.
func (engine *RPCClientEngine) Parse(text []byte) (Template, error) {
	return &remoteTemplate{engine, text}, nil
}

func (engine *RPCClientEngine) createSession() error {
	res, err := engine.rpcClient.Call("session")
	if err != nil {
		return err
	}
	if res.Error != nil {
		return &RPCError{*res.Error}
	}
	var result struct {
		SessionID  string
		RPCVersion string
	}
	if err := res.GetObject(&result); err != nil {
		return err
	}
	if result.RPCVersion != RPCVersion {
		return fmt.Errorf("Liquid server RPC mismatch: expected %s; actual %s", RPCVersion, result.RPCVersion)
	}
	engine.rpcClient.SetCustomHeader("Session-ID", result.SessionID)
	return err
}

func (engine *RPCClientEngine) rpcCall(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	res, err := engine.rpcClient.Call(method, params...)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, &RPCError{*res.Error}
	}
	return res, nil
}

// FileURLMap sets the filename -> permalink map that is used during link tag expansion.
func (engine *RPCClientEngine) FileURLMap(m map[string]string) (err error) {
	_, err = engine.rpcCall("fileUrls", m)
	return
}

// IncludeDirs specifies the search directories for the include tag.
func (engine *RPCClientEngine) IncludeDirs(dirs []string) (err error) {
	abs := make([]string, len(dirs))
	for i, dir := range dirs {
		abs[i], err = filepath.Abs(dir)
		if err != nil {
			return
		}
	}
	_, err = engine.rpcCall("includeDirs", abs)
	return
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
