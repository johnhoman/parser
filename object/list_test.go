package object

import (
    "github.com/stretchr/testify/require"
    "testing"
)

func TestList_Inspect(t *testing.T) {
    list := &List{Values: []Object{
        &Integer{Value: 1},
        &Integer{Value: 2},
        &Integer{Value: 3},
        &Integer{Value: 1},
        &Boolean{Value: true},
        &String{Value: "the quick brown fox"},
    }}

    require.Equal(t, `[1, 2, 3, 1, true, "the quick brown fox"]`, list.Inspect())
}

func TestList_Len(t *testing.T) {
    list := &List{Values: []Object{
        &Integer{Value: 1},
        &Integer{Value: 2},
        &Integer{Value: 3},
        &Integer{Value: 1},
        &Boolean{Value: true},
        &String{Value: "the quick brown fox"},
    }}
    require.Equal(t, &Integer{Value: 6}, list.Len())
}
