package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dadosjusbr/coletores/status"
	"github.com/dadosjusbr/datapackage"
	"github.com/dadosjusbr/proto/pipeline"

	"google.golang.org/protobuf/encoding/prototext"
)

func main() {
	// Reading input.
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		err := status.NewError(status.InvalidInput, fmt.Errorf("error reading execution result from stdin: %q", err))
		status.ExitFromError(err)
	}
	var er pipeline.ResultadoExecucao
	if err := prototext.Unmarshal(in, &er); err != nil {
		err := status.NewError(status.InvalidInput, fmt.Errorf("error unmarshalling execution result: %q", err))
		status.ExitFromError(err)
	}

	// Loading and validating package.xs
	if _, err := datapackage.Load(er.Pr.Pacote); err != nil {
		err = status.NewError(status.InvalidInput, fmt.Errorf("error loading datapackage (%s):%q", er.Pr.Pacote, err))
		status.ExitFromError(err)
	}

	// Printing output.
	b, err := prototext.Marshal(&er)
	if err != nil {
		err = status.NewError(status.Unknown, fmt.Errorf("error marshalling output:%q", err))
		status.ExitFromError(err)
	}
	fmt.Printf("%s", b)
}
