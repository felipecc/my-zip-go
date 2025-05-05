# Meu compactador de arquivos


## A ideia

A ideia geral do projeto seria entender como um compactor funciona, e como ele pode ser feito em Go. Inicialmente usei o zlib como lib de compressão, e o restante da estrutura da implementação foi feita para que eu entenda o como fazer isso em `go` e ainda conseguisse uma organização da estrutura.
O que imaginei, foi criar uma estrutura de saída de arquivo  que fosse divida em duas partes:

1. Um cabecalho que contenha o nome do arquivo, seu tamanho original e o tamanho comprimido.
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
 *--------
 | NameLength |
 *--------
 | OriginalSize |
 *--------
 | CompressedSize |
 *--------
```

```go
type  struct {
	NameLength     uint32
	OriginalSize   uint32
	CompressedSize uint32
}
```

### CompressBytes

A função `CompressBytes` recebe um slice de bytes e retorna um slice de bytes comprimidos ou um erro.

```go
func CompressBytes(data []byte) ([]byte, error)
```


