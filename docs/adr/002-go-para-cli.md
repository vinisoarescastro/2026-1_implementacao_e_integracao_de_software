# ADR-002 — Go como linguagem dos CLIs

## Contexto

O projeto exige dois CLIs multiplataforma (`assinatura` e `simulador`) que funcionem em Windows, Linux e macOS sem que o usuário precise instalar um runtime adicional. As alternativas consideradas foram:

- **Go**: cross-compilation nativa, binário estático único, biblioteca padrão rica para HTTP, processos e I/O.
- **Python**: requer runtime instalado; empacotamento multiplataforma (PyInstaller) é frágil.
- **Node.js**: requer runtime instalado; binários grandes com pkg/nexe.
- **Rust**: cross-compilation viável, mas curva de aprendizado mais alta e ecossistema de CLI menos maduro para este contexto.

A especificação (`plano-revisitado-v2.md`) já define Go 1.25 como premissa.

## Decisão

Os CLIs são desenvolvidos em **Go 1.25** usando a biblioteca [Cobra](https://github.com/spf13/cobra) para parsing de comandos.

## Consequências

**Fica mais fácil:**
- Cross-compilation para `windows/amd64`, `linux/amd64` e `darwin/amd64` com um único comando (`GOOS=... GOARCH=... go build`).
- Binário estático sem dependências — o usuário baixa e executa, sem instalar Go.
- Biblioteca padrão cobre HTTP client, execução de subprocessos e manipulação de arquivos sem dependências externas.
- Cobra gera `--help` estruturado e suporte a subcomandos com pouco código.

**Fica mais difícil:**
- A equipe precisa conhecer Go; não é a linguagem principal do curso.
- Gerenciamento de goroutines requer atenção para não introduzir race conditions nos modos de detecção de servidor.
