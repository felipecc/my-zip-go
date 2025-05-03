package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type FileHeader struct {
	NameLength     uint32
	OriginalSize   uint32
	CompressedSize uint32
}

func CompressBytes(data []byte) ([]byte, error) {
	var buff bytes.Buffer

	zNew := zlib.NewWriter(&buff)

	_, err := zNew.Write(data)

	if err != nil {
		return nil, fmt.Errorf("error compressing data: %v", err)
	}

	err = zNew.Close()

	if err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	return buff.Bytes(), nil
}

func DecompressBytes(data []byte) ([]byte, error) {

	rNew, err := zlib.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %v", err)
	}

	defer rNew.Close()

	var out bytes.Buffer

	_, err = io.Copy(&out, rNew)

	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %v", err)
	}

	return out.Bytes(), nil
}

func writeCompressedFile(fileName string, compressedData []byte) error {
	file, err := os.Stat(fileName)
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	originalFileSize := uint32(file.Size())

	fileNameInBytes := []byte(file.Name())
	fileNameLengh := uint32(len(fileNameInBytes))

	outPutFile, err := os.Create(fmt.Sprintf("%s.%s", fileName, "myz"))

	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}

	defer outPutFile.Close()

	header := FileHeader{
		NameLength:     fileNameLengh,
		OriginalSize:   originalFileSize,
		CompressedSize: uint32(len(compressedData)),
	}

	// escreve o header no arquivo como binario, aqui a forma de escrever é diferente
	binary.Write(outPutFile, binary.LittleEndian, header)

	// escreve o nome do arquivo em []byte
	outPutFile.Write(fileNameInBytes)

	// escreve o conteudo comprimido em []byte
	outPutFile.Write(compressedData)

	return nil
}

func readCompressedFile(fileWithPath string, outPutFileName string) error {

	// abre o arquivo para leitura
	file, err := os.Open(fileWithPath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var header FileHeader
	binary.Read(file, binary.LittleEndian, &header)

	fileNameInBytes := make([]byte, header.NameLength)
	n, err := file.Read(fileNameInBytes)
	if err != nil || uint32(n) != header.NameLength {
		return fmt.Errorf("error reading file name: %v", err)
	}
	fileName := string(fileNameInBytes)
	fmt.Println("fileName", fileName)

	// compressedData vai conter o conteudo comprimido que esta em bytes
	// ele aqui o tamanho vem do header.CompressedSize
	// o Read vai ler o conteudo do tamanho do header.CompressedSize
	compressedData := make([]byte, header.CompressedSize)
	n, err = file.Read(compressedData)
	if err != nil || uint32(n) != header.CompressedSize {
		return fmt.Errorf("error reading compressed data: %v", err)
	}

	decompressedData, err := DecompressBytes(compressedData)

	if err != nil {
		return fmt.Errorf("error decompressing data: %v", err)
	}

	outPutFile, err := os.Create(outPutFileName)

	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outPutFile.Close()

	n, err = outPutFile.Write(decompressedData)
	if err != nil || uint32(n) != header.OriginalSize {
		return fmt.Errorf("error writing decompressed data: %v", err)
	}

	return nil
}

func compressAndSave(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	compressedData, err := CompressBytes(file)
	if err != nil {
		return fmt.Errorf("error compressing data: %v", err)
	}

	err = writeCompressedFile(filePath, compressedData)
	if err != nil {
		return fmt.Errorf("error writing compressed file: %v", err)
	}

	return nil
}

func main() {
	// Primeiro comprimimos o arquivo
	err := compressAndSave("example.txt")
	if err != nil {
		fmt.Printf("Error in compressing file: %v\n", err)
		return
	}
	fmt.Println("File compressed successfully!")

	// Depois descomprimimos o arquivo
	err = readCompressedFile("example.txt.myz", "example_decompressed.txt")
	if err != nil {
		fmt.Printf("Error in decompressing file: %v\n", err)
		return
	}
	fmt.Println("File decompressed successfully!")
}
