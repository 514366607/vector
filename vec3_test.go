package vector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVec3(t *testing.T) {
	require.True(t, Vector3(0, 0, 0).Equals(Zero()))
	require.True(t, Vector3(1, 1, 1).Equals(One()))

	require.True(t, Vector3(0, 1, 0).Equals(UP()))
	require.True(t, Vector3(0, -1, 0).Equals(DOWN()))
	require.True(t, Vector3(-1, 0, 0).Equals(LEFT()))
	require.True(t, Vector3(1, 0, 0).Equals(RIGHT()))
	require.True(t, Vector3(0, 0, 1).Equals(FORWARD()))
	require.True(t, Vector3(0, 0, -1).Equals(BACK()))
}

func TestSlerp(t *testing.T) {
	a := Zero()
	b := LEFT()
	c := Slerp(a, b, 0.01)
	// t.Fatal(c)
	require.NotZero(t, c.X)
	require.Zero(t, c.Y)
	require.Zero(t, c.Z)
}

func TestLerp(t *testing.T) {
	a := Zero()
	b := LEFT()
	c := Lerp(a, b, 0.01)
	// t.Fatal(c)
	require.NotZero(t, c.X)
	require.Zero(t, c.Y)
	require.Zero(t, c.Z)
}

func TestMoveTowards(t *testing.T) {
	a := Zero()
	b := Vector3(100, 100, 0)
	c := MoveTowards(a, b, 10)
	require.Less(t, c.X, float64(10))
	require.Less(t, c.Y, float64(10))
	require.NotZero(t, c.X)
	require.NotZero(t, c.Y)
	require.Zero(t, c.Z)

	c = MoveTowards(a, b, Distance(a, b))
	require.True(t, c.Equals(b))

	require.Equal(t, MoveTowards(a, b, 0), Zero())
}

func TestMovePath(t *testing.T) {
	a := Zero()
	b := Vector3(100, 100, 0)

	c := []Vec3{}
	MovePath(a, b, 1, &c)
	require.Len(t, c, 143)

	c = c[:0]
	MovePath(Zero(), Vector3(0, 100, 0), 1, &c)
	require.Len(t, c, 101)
	// t.Fatal(MovePath(Zero(), Vector3(0, 100, 0), 1))
}

// cpu: Intel(R) Core(TM) i5-8500 CPU @ 3.00GHz
// BenchmarkMovePath-6   	  464095	      2521 ns/op	       0 B/op	       0 allocs/op
func BenchmarkMovePath(b *testing.B) {
	current := Zero()
	target := Vector3(1000, 1000, 0)
	c := []Vec3{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c = c[:0]
		MovePath(current, target, 1, &c)
	}
}
