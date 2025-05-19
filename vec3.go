package vector

import (
	"math"
)

// Vec3 Vec3
type Vec3 struct {
	X, Y, Z float64
}

var zeroVector = Vec3{0, 0, 0}

// Zero Zero
func Zero() Vec3 { return zeroVector }

var oneVector = Vec3{1, 1, 1}

// One One
func One() Vec3 { return oneVector }

var upVector = Vec3{0, 1, 0}

// UP UP
func UP() Vec3 { return upVector }

var downVector = Vec3{0, -1, 0}

// DOWN DOWN
func DOWN() Vec3 { return downVector }

var leftVector = Vec3{-1, 0, 0}

// LEFT LEFT
func LEFT() Vec3 { return leftVector }

var rightVector = Vec3{1, 0, 0}

// RIGHT RIGHT
func RIGHT() Vec3 { return rightVector }

var forwardVector = Vec3{0, 0, 1}

// FORWARD FORWARD
func FORWARD() Vec3 { return forwardVector }

var backVector = Vec3{0, 0, -1}

// BACK BACK
func BACK() Vec3 { return backVector }

// Slerp 实现球面线性插值
func Slerp(a, b Vec3, t float64) Vec3 {
	// 向量归一化处理
	aNorm := a.Normalized()
	bNorm := b.Normalized()

	// 计算点积和夹角
	dot := aNorm.Dot(bNorm)
	dot = float64(math.Max(-1, math.Min(1, float64(dot))))
	theta := float64(math.Acos(float64(dot))) * t

	// 计算插值向量
	relativeVec := bNorm.Sub(aNorm.Mul(dot)).Normalized()
	return aNorm.Mul(float64(math.Cos(float64(theta)))).Add(
		relativeVec.Mul(float64(math.Sin(float64(theta)))))
}

// Mul 向量标量乘法
func (v Vec3) Mul(scalar float64) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

// Add 向量加法
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Sub 向量减法
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Clamp Clamp
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Clamp01 Clamp01
func Clamp01(value float64) float64 {
	return Clamp(value, 0, 1)
}

// Lerp 在两个点之间进行线性插值
func Lerp(a, b Vec3, t float64) Vec3 {
	t = Clamp01(t)
	return Vec3{a.X + (b.X-a.X)*t, a.Y + (b.Y-a.Y)*t, a.Z + (b.Z-a.Z)*t}
}

// MoveTowards 计算出由当前点和目标点所确定的两点之间的位置，移动的距离不得超过由 maxDistanceDelta 参数所指定的值。
func MoveTowards(current, target Vec3, maxDistanceDelta float64) Vec3 {
	if maxDistanceDelta == 0 {
		return Zero()
	}

	num1 := target.X - current.X
	num2 := target.Y - current.Y
	num3 := target.Z - current.Z
	d := (num1*num1 + num2*num2 + num3*num3)
	if d == 0 || maxDistanceDelta >= 0 && d <= maxDistanceDelta*maxDistanceDelta {
		return target
	}
	num4 := math.Sqrt(d)
	return Vec3{current.X + num1/num4*maxDistanceDelta, current.Y + num2/num4*maxDistanceDelta, current.Z + num3/num4*maxDistanceDelta}
}

// MovePath 计算出由当前点和目标点所确定的两点之间经过的点，step 参数为每步距离
func MovePath(current, target Vec3, step float64, path *[]Vec3) {
	num1 := target.X - current.X
	num2 := target.Y - current.Y
	num3 := target.Z - current.Z
	d := (num1*num1 + num2*num2 + num3*num3)
	if d == 0 || step == 0 || step > 0 && d <= step*step {
		*path = append(*path, current, target)
		return
	}
	distance := math.Sqrt(d)

	*path = append(*path, current)

	num1 = num1 / distance
	num2 = num2 / distance
	num3 = num3 / distance

	for maxDistanceDelta := float64(step); maxDistanceDelta < distance; maxDistanceDelta += step {
		*path = append(*path, Vec3{current.X + num1*maxDistanceDelta, current.Y + num2*maxDistanceDelta, current.Z + num3*maxDistanceDelta})
	}
	*path = append(*path, target)
}

// SmoothDamp 平滑的阻尼效果
func SmoothDamp(current, target Vec3, currentVelocity *Vec3, smoothTime, deltaTime, maxSpeed float64) Vec3 {
	smoothTime = max(0.0001, smoothTime)
	num1 := 2 / smoothTime
	num2 := num1 * deltaTime
	num3 := (1.0 / (1.0 + num2 + 0.47999998927116394*num2*num2 + 0.23499999940395355*num2*num2*num2))
	num4 := current.X - target.X
	num5 := current.Y - target.Y
	num6 := current.Z - target.Z

	vector3 := target
	num7 := maxSpeed * smoothTime
	num8 := num7 * num7
	d := (num4*num4 + num5*num5 + num6*num6)
	if d > num8 {
		num9 := math.Sqrt(d)
		num4 = num4 / num9 * num7
		num5 = num5 / num9 * num7
		num6 = num6 / num9 * num7
	}

	target.X = current.X - num4
	target.Y = current.Y - num5
	target.Z = current.Z - num6
	num10 := (currentVelocity.X + num1*num4) * deltaTime
	num11 := (currentVelocity.Y + num1*num5) * deltaTime
	num12 := (currentVelocity.Z + num1*num6) * deltaTime
	currentVelocity.X = (currentVelocity.X - num1*num10) * num3
	currentVelocity.Y = (currentVelocity.Y - num1*num11) * num3
	currentVelocity.Z = (currentVelocity.Z - num1*num12) * num3
	x := target.X + (num4+num10)*num3
	y := target.Y + (num5+num11)*num3
	z := target.Z + (num6+num12)*num3
	num13 := vector3.X - current.X
	num14 := vector3.Y - current.Y
	num15 := vector3.Z - current.Z
	num16 := x - vector3.X
	num17 := y - vector3.Y
	num18 := z - vector3.Z
	if num13*num16+num14*num17+num15*num18 > 0 {
		x = vector3.X
		y = vector3.Y
		z = vector3.Z
		currentVelocity.X = (x - vector3.X) / deltaTime
		currentVelocity.Y = (y - vector3.Y) / deltaTime
		currentVelocity.Z = (z - vector3.Z) / deltaTime
	}
	return Vec3{x, y, z}
}

// Vector3 Vector3
func Vector3(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

// Scale 对两个向量的各分量进行相乘操作
func Scale(a, b Vec3) Vec3 {
	return Vec3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Cross 两个向量的乘积
func (v Vec3) Cross(a, b Vec3) Vec3 {
	return Vec3{(a.Y*b.Z - a.Z*b.Y), (a.Z*b.X - a.X*b.Z), (a.X*b.Y - a.Y*b.X)}
}

func (v Vec3) Equals(other Vec3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

func Equals(a, b Vec3) bool {
	return a.X == b.X && a.Y == b.Y && a.Z == b.Z
}

// // Reflects 垂直于由法线定义的平面的向量
func (v Vec3) Reflects(inNormal Vec3) Vec3 {
	num := -2 * inNormal.Dot(v)
	return Vec3{num*inNormal.X + v.X, num*inNormal.Y + v.Y, num*inNormal.Z + v.Z}
}

// Normalized 向量归一化
func (v Vec3) Normalized() Vec3 {
	magnitude := v.Magnitude()
	if magnitude > 9.999999747378752e-06 {
		return Vec3{v.X / magnitude, v.Y / magnitude, v.Z / magnitude}
	}
	return zeroVector
}

// Dot 向量点积计算
func (v Vec3) Dot(other Vec3) float64 {
	return Dot(v, other)
}

// Dot 向量点积计算
func Dot(a, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

// Project 将一个向量投影到另一个向量上
func Project(vector, onNormal Vec3) Vec3 {
	num1 := Dot(onNormal, onNormal)
	if num1 < 1.1754944e-38 {
		return zeroVector
	}
	num2 := Dot(vector, onNormal)
	return Vec3{onNormal.X * num2 / num1, onNormal.Y * num2 / num1, onNormal.Z * num2 / num1}
}

// Angle 返回从点和到点之间所形成的角度（以度为单位）
func Angle(from, to Vec3) float64 {
	num := math.Sqrt(from.SqrMagnitude() * to.SqrMagnitude())
	if num < 1.0000000036274937e-15 {
		return 0
	}
	return math.Acos(Clamp(Dot(from, to)/num, -1, 1)) * 57.29578
}

// Distance 返回两个点的距离
func Distance(a, b Vec3) float64 {
	num1 := a.X - b.X
	num2 := a.Y - b.Y
	num3 := a.Z - b.Z
	return math.Sqrt(num1*num1 + num2*num2 + num3*num3)
}

// ClampMagnitude 返回一个与原向量相同的向量副本，但其大小会被限制在最大长度范围内。
func ClampMagnitude(a, b Vec3, maxLength float64) Vec3 {
	sqrMagnitude := a.SqrMagnitude()
	if sqrMagnitude <= maxLength*maxLength {
		return a
	}
	num1 := math.Sqrt(sqrMagnitude)
	num2 := a.X / num1
	num3 := a.Y / num1
	num4 := a.Z / num1
	return Vec3{num2 * maxLength, num3 * maxLength, num4 * maxLength}
}

// Magnitude Magnitude
func (v Vec3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// SqrMagnitude 返回此向量的长度
func SqrMagnitude(vector Vec3) float64 {
	return vector.SqrMagnitude()
}

// SqrMagnitude 返回此向量的长度
func (v Vec3) SqrMagnitude() float64 {
	return (v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Min 返回一个由两个向量中最小的元素构成的向量
func Min(a, b Vec3) Vec3 {
	return Vec3{min(a.X, b.X), min(a.Y, b.Y), min(a.Z, b.Z)}
}

// Max 返回一个由两个向量中最大值构成的向量
func Max(a, b Vec3) Vec3 {
	return Vec3{max(a.X, b.X), max(a.Y, b.Y), max(a.Z, b.Z)}
}
