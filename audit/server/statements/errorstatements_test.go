package statements

import (
	"cdcon21builddriver/globalUtils"
	"fmt"
	"testing"
)

func Test_UserErr_internalError(t *testing.T) {
	var usrerr UserErr
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		ge      UserErr
		args    args
		want    string
		lang    globalUtils.Languages
		setLang bool
	}{
		{name: "Pull English internal error", ge: usrerr, args: args{err: fmt.Errorf("test Error")}, want: "Internal error. Error: test Error\n",
			lang: globalUtils.LangEN, setLang: false,
		},
		{name: "Pull English internal error", ge: usrerr, args: args{err: fmt.Errorf("error prueba")}, want: "Error interno. Error: error prueba\n",
			lang: globalUtils.LangES, setLang: true,
		},
		{name: "Pull English internal error", ge: usrerr, args: args{err: fmt.Errorf("test Error")}, want: "Internal error. Error: test Error\n",
			lang: globalUtils.LangEN, setLang: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setLang {
				SetLanguage(tt.lang)
			}
			if got := tt.ge.internalError(tt.args.err); got != tt.want {
				t.Errorf("InsertError() = %v, want %v for language %v", got, tt.want, tt.lang)
			}
		})
	}
}

func Test_UserErr_getSqlTxt(t *testing.T) {
	var usrerr UserErr
	type args struct {
		errKey   string
		language globalUtils.Languages
	}
	tests := []struct {
		name string
		ge   UserErr
		args args
		want string
	}{
		{name: "Pull English internal error", ge: usrerr, args: args{errKey: "internalError", language: globalUtils.LangEN}, want: "Internal error. Error: %v\n"},
		{name: "Pull Spanish internal error", ge: usrerr, args: args{errKey: "internalError", language: globalUtils.LangES}, want: "Error interno. Error: %v\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ge.getSqlTxt(tt.args.errKey, tt.args.language); got != tt.want {
				t.Errorf("getSqlTxt() = %v, want %v", got, tt.want)
			}
		})
	}
}
