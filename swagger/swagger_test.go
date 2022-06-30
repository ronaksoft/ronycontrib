package swagger_test

import (
    "encoding/json"
    "fmt"
    "strings"
    "testing"

    "github.com/clubpay/ronycontrib/swagger"
    "github.com/clubpay/ronykit/desc"
    "github.com/clubpay/ronykit/std/gateway/fasthttp"
)

type sampleReq struct {
    X string   `json:"x"`
    Y string   `json:"y"`
    Z int64    `json:"z"`
    W []string `json:"w"`
}

type subRes struct {
    Some    string `json:"some"`
    Another []byte `json:"another"`
}

type sampleRes struct {
    Out1 int      `json:"out1"`
    Out2 string   `json:"out2"`
    Sub  subRes   `json:"sub"`
    Subs []subRes `json:"subs"`
}

type sampleError struct {
    Code        int    `json:"code" swag:"enum:504,503"`
    Description string `json:"description"`
}

type anotherRes struct {
    subRes
    Out1 int    `json:"out1"`
    Out2 string `json:"out2"`
}

type testService struct{}

func (t testService) Desc() *desc.Service {
    return (&desc.Service{
        Name:         "testService",
        PreHandlers:  nil,
        PostHandlers: nil,
    }).
            AddContract(
                desc.NewContract().
                        AddSelector(fasthttp.Selector{
                            Method: fasthttp.MethodGet,
                            Path:   "/some/:x/:y",
                        }).
                    SetInput(&sampleReq{}).
                    SetOutput(&sampleRes{}).
                    AddPossibleError(404, "ITEM1", &sampleError{}).
                    AddPossibleError(504, "SERVER", &sampleError{}).
                    SetHandler(nil),
            ).
            AddContract(
                desc.NewContract().
                        AddSelector(fasthttp.Selector{
                            Method: fasthttp.MethodPost,
                            Path:   "/some/:x/:y",
                        }).
                    SetInput(&sampleReq{}).
                    SetOutput(&anotherRes{}).
                    SetHandler(nil),
            )
}

func TestNewSwagger(t *testing.T) {
    sg := swagger.NewSwagger("TestTitle", "v0.0.1", "")
    sg.WithTag("json")

    sb := &strings.Builder{}
    err := sg.WriteTo(sb, testService{})
    if err != nil {
        t.Fatal(err)
    }

    x, _ := json.MarshalIndent(json.RawMessage(sb.String()), "", "   ")
    fmt.Println(string(x))
}
