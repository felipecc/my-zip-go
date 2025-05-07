# Meu compactador de arquivos


## A ideia

A ideia geral do projeto seria entender como um compactor funciona, e como ele pode ser feito em Go. Inicialmente usei o zlib como lib de compressão, e o restante da estrutura da implementação foi feita para que eu entenda o como fazer isso em `go` e ainda conseguisse uma organização da estrutura.
O que imaginei, foi criar uma estrutura de saída de arquivo  que fosse divida em duas partes:

1. Um cabeçalho que contenha o nome do arquivo, seu tamanho original e o tamanho comprimido.
2. O conteudo do arquivo.

Como programa  seria necessário criar as seguintes funções:

1. Uma função que dado a entrada em bytes, retorne o conteudo bytes comprimidos.
2. Uma função que dado a entrada em bytes comprimidos, retorne o conteudo bytes descomprimidos.

3. Uma função que recebesse um nome de arquivo como parametro e um conteudo em bytes comprimidos tambem como parametro, coordena a escrita do conteudo em bytes comprimidos no arquivo.


4. Uma função que recebesse um path com um arquivo comprimido como parametro e um nome alternativo de saida também como parametro, para coordenar a leitura, descomprimir o arquivo e escrever o conteudo descomprimido em um novo arquivo com o nome alternativo de saida.


## Implementação
O arquivo de saída terá o diagrama de estrutura de arquivo a seguir:

```
 *----------------------------------*
 | FileHeader |  CompressedData     |
 *----------------------------------*
```

### FileHeader

FileHeader é a estrutura de arquivo que contem o nome do arquivo, seu tamanho original e o tamanho comprimido. No arquivo de saída vamos escreve-la em binário, para eliminar a etapa da tradução da estrutura.

```
 *----------------*
 | NameLength     |
 *----------------*
 | OriginalSize   |
 *----------------*
 | CompressedSize |
 *----------------*
```

```go
type  struct {
	NameLength     uint32
	OriginalSize   uint32
	CompressedSize uint32
}
```

### CompressBytes

A função `CompressBytes` recebe um slice de bytes e retorna um slice de bytes para serem comprimidos ou um erro.

```go
func CompressBytes(data []byte) ([]byte, error) {
	
	var buff bytes.Buffer

	zw := zlib.NewWriter(&buff)

	_, err := zw.Write(data)

	if err != nil {
		return nil, fmt.Errorf("error compressing data: %v", err)
	}

	err = zw.Close()

	if err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	return buff.Bytes(), nil
}
```

Na função `CompressBytes`, o processo de compressão ocorre da seguinte forma:

1. Primeiro, criamos um buffer (`buff`) usando `bytes.Buffer` para armazenar os dados comprimidos
2. Em seguida, criamos um escritor zlib (`zw`) usando `zlib.NewWriter(&buff)`, que implementa a interface `io.Writer`
3. Utilizamos o método `Write` do escritor zlib para escrever os dados (`data`) no buffer, realizando a compressão
4. Fechamos o escritor zlib usando `zw.Close()` para garantir que todos os dados sejam escritos no buffer
5. Por fim, retornamos os bytes comprimidos usando `buff.Bytes()`

### DecompressBytes

A função `DecompressBytes` recebe um slice de bytes compromidos e retorna um slice de bytes ou um erro.

```go
func DecompressBytes(data []byte) ([]byte, error) {

	zr, err := zlib.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %v", err)
	}

	defer zr.Close()

	var out bytes.Buffer

	_, err = io.Copy(&out, zr)

	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %v", err)
	}

	return out.Bytes(), nil
}
```

Na função `DecompressBytes`, o processo de descompressão ocorre da seguinte forma:

1. Convertemos os dados comprimidos (`data`) em um `bytes.Reader` usando `bytes.NewReader(data)`
2. Em seguida, criamos um leitor zlib (`zr`) usando `zlib.NewReader`, que implementa a interface `io.Reader`
3. Criado um buffer (`out`) para armazenar os dados descomprimidos
4. O `io.Copy` para copiar os dados do leitor zlib (`zr`) para o buffer (`out`). A ordem dos parâmetros é importante: primeiro o destino (`&out`), depois a fonte (`zr`)
5. Retornamos os bytes descomprimidos usando `out.Bytes()`

O uso de `&out` é necessário porque o `io.Copy` precisa modificar o buffer para armazenar os dados descomprimidos, então passamos o endereço de memória do buffer.


### writeCompressedFile

A função `writeCompressedFile` recebe um nome de arquivo e um slice de bytes comprimidos e escreve o conteúdo em um arquivo com o nome de saída.

```go
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
```

os.Stat retorna uma implementação da interface `fs.FileInfo` que contem informações sobre o arquivo.

```go
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() any           // underlying data source (can return nil)
}
```
Estas informações são usadas para preencher o `FileHeader`. Com o `os.Create` criamos o arquivo de saída e o retornamos como `outPutFile`.
O `defer outPutFile.Close()` é usado para garantir que o arquivo seja fechado quando a função `writeCompressedFile` terminar.
Para escrever o cabeçalho, vamos escrever em binário usando `binary.Write`, nesta etapa economizamos a etapa de serialização do cabeçalho. Os parametros da função `binary.Write` são:













