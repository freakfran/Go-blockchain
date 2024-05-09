package core

import "io"

// Encoder 接口定义了编码器的行为，允许将数据编码并写入指定的 writer 中。
// [T any] 使用了泛型，使得 Encoder 可以支持任意类型的编码。
type Encoder[T any] interface {
	// Encode 方法将指定的数据 data 编码，并写入到提供的 writer 中。
	// writer: 用于写入编码后数据的目标 io.Writer。
	// data: 需要被编码的数据，其类型为泛型 T。
	// 返回值 error: 如果编码过程中发生错误，则返回非 nil 的 error 对象。
	Encode(writer io.Writer, data T) error
}

// Decoder 接口定义了解码器的行为，允许从指定的 reader 中读取数据并解码到指定的对象中。
// [T any] 使用了泛型，使得 Decoder 可以支持任意类型的解码。
type Decoder[T any] interface {
	// Decode 方法从指定的 reader 中读取数据，并将其解码到提供的 data 对象中。
	// reader: 用于读取解码数据的来源 io.Reader。
	// data: 用于存储解码后数据的对象，其类型为泛型 T。
	// 返回值 error: 如果解码过程中发生错误，则返回非 nil 的 error 对象。
	Decode(reader io.Reader, data T) error
}
