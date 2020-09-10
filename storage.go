package jda

import (
	"os"
	"reflect"
	"fmt"
	"unsafe"
	"sync"
)

func StorageGetStructOnPosition(
	f *os.File,
	fileMutex *sync.Mutex,
	position int64,
	outputInter interface{},
) error {
	l := GetLogger()

	fileMutex.Lock()

	_, err := f.Seek(position, os.SEEK_SET)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to change file pointer")
		fileMutex.Unlock()
		return l.ErrorQueue
	}

	structValue := reflect.Indirect(reflect.ValueOf(outputInter))
	length := structValue.NumField()
	for i := 0; i < length; i = i + 1 {
		value := structValue.Field(i)
		fieldInter := value.Interface()
		switch ty := fieldInter.(type) {
		case int:
			var number int
			_, err := f.Read((*(*[unsafe.Sizeof(number)]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read int in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(number))
		case int32:
			var number int32
			_, err := f.Read((*(*[4]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read int32 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(number))
		case int64:
			var number int64
			_, err := f.Read((*(*[8]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read int64 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(number))
		case uint:
			var number uint
			_, err := f.Read((*(*[unsafe.Sizeof(number)]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read uint in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(number))
		case uint32:
			var number uint32
			_, err := f.Read((*(*[4]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read uint32 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(number))
		case uint64:
			var number uint64
			_, err := f.Read((*(*[8]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read uint64 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(number))
		case []byte:
			dataLength := uint64(0)
			_, err := f.Read((*(*[8]byte)(unsafe.Pointer(&dataLength)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read []byte field data size in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			buffer := make([]byte, dataLength)
			_, err = f.Read(buffer[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read []byte in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(buffer))
		case string:
			dataLength := uint64(0)
			_, err := f.Read((*(*[8]byte)(unsafe.Pointer(&dataLength)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read string field data size in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			buffer := make([]byte, dataLength)
			_, err = f.Read(buffer[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in read string in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			value.Set(reflect.ValueOf(string(buffer)))
		default:
			l.Error("Invalid type "+fmt.Sprint(ty))
			fileMutex.Unlock()
			return l.ErrorQueue
		}
	}

	fileMutex.Unlock()
	return nil
}

func StorageInsertStructOnPosition(
	f *os.File,
	fileMutex *sync.Mutex,
	position int64,
	inter interface{},
) error {
	l := GetLogger()

	fileMutex.Lock()

	_, err := f.Seek(position, os.SEEK_SET)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to change file pointer")
		fileMutex.Unlock()
		return l.ErrorQueue
	}

	structValue := reflect.ValueOf(inter)
	length := structValue.NumField()
	for i := 0; i < length; i = i + 1 {
		fieldInter := structValue.Field(i).Interface()
		switch ty := fieldInter.(type) {
		case int:
			number := fieldInter.(int)
			_, err := f.Write((*(*[unsafe.Sizeof(number)]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write int in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case int32:
			number := fieldInter.(int32)
			_, err := f.Write((*(*[4]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write int32 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case int64:
			number := fieldInter.(int64)
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write int64 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case uint:
			number := fieldInter.(uint)
			_, err := f.Write((*(*[unsafe.Sizeof(number)]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write uint in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case uint32:
			number := fieldInter.(uint32)
			_, err := f.Write((*(*[4]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write uint32 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case uint64:
			number := fieldInter.(uint64)
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write uint64 in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case []byte:
			data := fieldInter.([]byte)
			dataLength := uint64(len(data))
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&dataLength)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write []byte field data size in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			_, err = f.Write(data)
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write []byte in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		case string:
			data := []byte(fieldInter.(string))
			dataLength := uint64(len(data))
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&dataLength)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write string field data size in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
			_, err = f.Write(data)
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write string in file")
				fileMutex.Unlock()
				return l.ErrorQueue
			}
		default:
			l.Error("Invalid type "+fmt.Sprint(ty))
			fileMutex.Unlock()
			return l.ErrorQueue
		}
	}

	fileMutex.Unlock()
	return nil
}

func StorageAppendStruct(
	f *os.File,
	fileMutex *sync.Mutex,
	inter interface{},
) (uint64, error) {
	l := GetLogger()

	fileMutex.Lock()

	size, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to change file pointer")
		fileMutex.Unlock()
		return 0, l.ErrorQueue
	}

	structValue := reflect.ValueOf(inter)
	length := structValue.NumField()
	for i := 0; i < length; i = i + 1 {
		fieldInter := structValue.Field(i).Interface()
		switch ty := fieldInter.(type) {
		case int:
			number := fieldInter.(int)
			_, err := f.Write((*(*[unsafe.Sizeof(number)]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write int in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case int32:
			number := fieldInter.(int32)
			_, err := f.Write((*(*[4]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write int32 in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case int64:
			number := fieldInter.(int64)
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write int64 in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case uint:
			number := fieldInter.(uint)
			_, err := f.Write((*(*[unsafe.Sizeof(number)]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write uint in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case uint32:
			number := fieldInter.(uint32)
			_, err := f.Write((*(*[4]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write uint32 in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case uint64:
			number := fieldInter.(uint64)
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&number)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write uint64 in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case []byte:
			data := fieldInter.([]byte)
			dataLength := uint64(len(data))
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&dataLength)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write []byte field data size in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
			_, err = f.Write(data)
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write []byte in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		case string:
			data := []byte(fieldInter.(string))
			dataLength := uint64(len(data))
			_, err := f.Write((*(*[8]byte)(unsafe.Pointer(&dataLength)))[:])
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write string field data size in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
			_, err = f.Write(data)
			if err != nil {
				l.Error(err.Error())
				l.Error("Error in write string in file")
				fileMutex.Unlock()
				return 0, l.ErrorQueue
			}
		default:
			l.Error("Invalid type "+fmt.Sprint(ty))
			fileMutex.Unlock()
			return 0, l.ErrorQueue
		}
	}

	fileMutex.Unlock()
	return uint64(size), nil
}