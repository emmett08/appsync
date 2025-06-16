
package port

type DirWriter interface {
    Write(path string, data []byte) error
}
