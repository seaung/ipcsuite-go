package camera

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/seaung/ipcsuite-go/internal/protos"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/seaung/ipcsuite-go/pkg/utils"
)

type CelOpt struct {
	envOptions     []cel.EnvOption
	programOptions []cel.ProgramOption
}

func NewEnvCelOption(c *CelOpt) (*cel.Env, error) {
	return cel.NewEnv(cel.Lib(c))
}

func (c *CelOpt) CompileOptions() []cel.EnvOption {
	return c.envOptions
}

func (c *CelOpt) ProgramOptions() []cel.ProgramOption {
	return c.programOptions
}

func (c *CelOpt) UpdateCompileOptions(params map[string]string) {
	for key, value := range params {
		var decl *exprpb.Decl
		if strings.HasPrefix(value, "randomInt") {
			decl = decls.NewIdent(key, decls.Int, nil)
		} else {
			decl = decls.NewIdent(key, decls.String, nil)
		}
		c.envOptions = append(c.envOptions, cel.Declarations(decl))
	}
}

func NewCelEnvOptions() CelOpt {
	celOptions := CelOpt{}

	celOptions.envOptions = []cel.EnvOption{
		cel.Container("protos"),
		cel.Types(
			&protos.UrlType{},
			&protos.Request{},
			&protos.Response{},
		),
		cel.Declarations(
			decls.NewIdent("request", decls.NewObjectType("protos.Request"), nil),
			decls.NewIdent("response", decls.NewObjectType("protos.Response"), nil),
		),
		cel.Declarations(
			decls.NewFunction("bcontains",
				decls.NewInstanceOverload("bytes_bcontains_bytes",
					[]*exprpb.Type{decls.Bytes, decls.Bytes},
					decls.Bool,
				)),
			decls.NewFunction("bmatches",
				decls.NewInstanceOverload("string_bmatches_bytes",
					[]*exprpb.Type{decls.String, decls.Bytes},
					decls.Bool,
				)),
			decls.NewFunction("md5",
				decls.NewInstanceOverload("md5_string",
					[]*exprpb.Type{decls.String},
					decls.String,
				)),
			decls.NewFunction("randomInt",
				decls.NewInstanceOverload("randomInt_int_int",
					[]*exprpb.Type{decls.Int, decls.Int},
					decls.Int,
				)),
			decls.NewFunction("randomLowercase",
				decls.NewInstanceOverload("randomLowercase_int",
					[]*exprpb.Type{decls.Int},
					decls.String,
				)),
			decls.NewFunction("base64",
				decls.NewInstanceOverload("base64_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String,
				)),
			decls.NewFunction("base64",
				decls.NewInstanceOverload("base64_string",
					[]*exprpb.Type{decls.String},
					decls.String,
				)),
			decls.NewFunction("base64Decode",
				decls.NewInstanceOverload("base64Decode_string",
					[]*exprpb.Type{decls.String},
					decls.String,
				)),
			decls.NewFunction("base64Decode",
				decls.NewInstanceOverload("base64Decode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String,
				)),
			decls.NewFunction("urlencode",
				decls.NewInstanceOverload("urlencode_string",
					[]*exprpb.Type{decls.String},
					decls.String,
				)),
			decls.NewFunction("urlencode",
				decls.NewInstanceOverload("urlencode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.Bytes,
				)),
			decls.NewFunction("urldecode",
				decls.NewInstanceOverload("urldecode_string",
					[]*exprpb.Type{decls.String},
					decls.String,
				)),
			decls.NewFunction("urldecode",
				decls.NewInstanceOverload("urldecode_bytes",
					[]*exprpb.Type{decls.Bytes},
					decls.String,
				)),
			decls.NewFunction("substr",
				decls.NewInstanceOverload("substr_string_int_int",
					[]*exprpb.Type{decls.String, decls.Int, decls.Int},
					decls.String,
				)),
			decls.NewFunction("icontains",
				decls.NewInstanceOverload("icontains_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool,
				)),
		),
	}

	celOptions.programOptions = []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "bytes_bcontains_bytes",
				Binary: func(lhs, rhs ref.Val) ref.Val {
					val1, ok := lhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bcontains", lhs.Type())
					}

					val2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bcontains", rhs.Type())
					}

					return types.Bool(bytes.Contains(val1, val2))
				},
			},
			&functions.Overload{
				Operator: "string_bmatches_bytes",
				Binary: func(lhs, rhs ref.Val) ref.Val {
					val1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bmatches", lhs.Type())
					}

					val2, ok := rhs.(types.Bytes)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bmatches", rhs.Type())
					}

					matcher, err := regexp.Match(string(val1), val2)
					if err != nil {
						return types.NewErr("%v", err)
					}

					return types.Bool(matcher)
				},
			},
			&functions.Overload{
				Operator: "md5_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to md5", value.Type())
					}

					return types.String(fmt.Sprintf("%x", md5.Sum([]byte(v))))
				},
			},
			&functions.Overload{
				Operator: "randomInt_int_int",
				Binary: func(lhs, rhs ref.Val) ref.Val {
					val1, ok := lhs.(types.Int)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to randomInt", lhs.Type())
					}

					val2, ok := rhs.(types.Int)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to randomInt", rhs.Type())
					}

					min, max := int(val1), int(val2)

					return types.Int(rand.Intn(max-min) + min)
				},
			},
			&functions.Overload{
				Operator: "randomLowercase_int",
				Unary: func(value ref.Val) ref.Val {
					number, ok := value.(types.Int)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to randomLowercase", value.Type())
					}

					return types.String(utils.RandomLowercase(int(number)))
				},
			},
			&functions.Overload{
				Operator: "base64_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_string", value.Type())
					}

					return types.String(base64.StdEncoding.EncodeToString([]byte(v)))
				},
			},
			&functions.Overload{
				Operator: "base64_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64_bytes", value.Type())
					}

					return types.String(base64.StdEncoding.EncodeToString(v))
				},
			},
			&functions.Overload{
				Operator: "base64Decode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_string", value.Type())
					}

					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}

					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "base64Decode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to base64Decode_bytes", value.Type())
					}

					decodeBytes, err := base64.StdEncoding.DecodeString(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}

					return types.String(decodeBytes)
				},
			},
			&functions.Overload{
				Operator: "urlencode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
					}

					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urlencode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urlencode_bytes", value.Type())
					}

					return types.String(url.QueryEscape(string(v)))
				},
			},
			&functions.Overload{
				Operator: "urldecode_string",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.String)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_string", value.Type())
					}

					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}

					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "urldecode_bytes",
				Unary: func(value ref.Val) ref.Val {
					v, ok := value.(types.Bytes)
					if !ok {
						return types.ValOrErr(value, "unexpected type '%v' passed to urldecode_bytes", value.Type())
					}

					decodeString, err := url.QueryUnescape(string(v))
					if err != nil {
						return types.NewErr("%v", err)
					}

					return types.String(decodeString)
				},
			},
			&functions.Overload{
				Operator: "substr_string_int_int",
				Function: func(values ...ref.Val) ref.Val {
					if len(values) == 3 {
						str, ok := values[0].(types.String)
						if !ok {
							return types.NewErr("invalid string to 'substr'")
						}

						start, ok := values[1].(types.Int)
						if !ok {
							return types.NewErr("invalid start to 'substr'")
						}

						length, ok := values[2].(types.Int)
						if !ok {
							return types.NewErr("invalid length to 'substr'")
						}

						runer := []rune(str)
						if start < 0 || length < 0 || int(start+length) > len(runer) {
							return types.NewErr("invalid start or length to 'substr'")
						}

						return types.String(runer[start : start+length])
					} else {
						return types.NewErr("too many arguments to 'substr'")
					}
				},
			},
			&functions.Overload{
				Operator: "icontains_string",
				Binary: func(lhs, rhs ref.Val) ref.Val {
					val1, ok := lhs.(types.String)
					if !ok {
						return types.ValOrErr(lhs, "unexpected type '%v' passed to bcontains", lhs.Type())
					}

					val2, ok := rhs.(types.String)
					if !ok {
						return types.ValOrErr(rhs, "unexpected type '%v' passed to bcontains", rhs.Type())
					}

					return types.Bool(strings.Contains(strings.ToLower(string(val1)), strings.ToLower(string(val2))))
				},
			},
		),
	}

	return celOptions
}

/*
执行表达式
*/
func EvalExpression(env *cel.Env, expression string, params map[string]interface{}) (ref.Val, error) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, iss.Err()
	}

	prag, err := env.Program(ast)
	if err != nil {
		return nil, err
	}

	out, _, err := prag.Eval(params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func UrlType2String(u *protos.UrlType) string {
	var buffer strings.Builder

	if u.Scheme != "" {
		buffer.WriteString(u.Scheme)
		buffer.WriteByte(':')
	}

	if u.Scheme != "" || u.Host != "" {
		if u.Host != "" || u.Path != "" {
			buffer.WriteString("//")
		}

		if h := u.Host; h != "" {
			buffer.WriteString(u.Host)
		}
	}

	path := u.Path

	if path != "" && path[0] != '/' && u.Host != "" {
		buffer.WriteByte('/')
	}

	if buffer.Len() == 0 {
		if i := strings.IndexByte(path, ':'); i > -1 && strings.IndexByte(path[:i], '/') == -1 {
			buffer.WriteString("./")

		}
	}

	buffer.WriteString(path)

	if u.Query != "" {
		buffer.WriteByte('?')
		buffer.WriteString(u.Query)
	}

	if u.Fragment != "" {
		buffer.WriteByte('#')
		buffer.WriteString(u.Fragment)
	}

	return buffer.String()
}
