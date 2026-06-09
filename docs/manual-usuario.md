# Manual do Usuário — Sistema Runner

O Sistema Runner é um conjunto de ferramentas de linha de comando (CLI) que permite criar e validar assinaturas digitais sem precisar conhecer os detalhes do ambiente Java.

## Pré-requisitos

| Ferramenta | Obrigatório | Como verificar |
|---|---|---|
| Go | Sim (v1.25+) | `go version` |
| Java JDK 21 | Não — baixado automaticamente se ausente | `java -version` |
| Maven | Não — apenas para desenvolvimento | `mvn -version` |

> O CLI detecta o JDK automaticamente na seguinte ordem: `~/.hubsaude/jdk/` → PATH do sistema → **download automático** do Eclipse Temurin 21.

---

## 1. Configuração inicial (primeira vez)

O CLI gerencia o `assinador.jar` e o JDK automaticamente. Na primeira execução, ele verifica a versão mais recente no repositório e faz o download se necessário.

Se preferir compilar o JAR localmente (modo desenvolvimento), siga os passos opcionais abaixo:

### 1.1. (Opcional) Compilar o assinador.jar localmente

Na raiz do projeto, execute:

```powershell
mvn -f assinador/pom.xml package -q
```

O comando gera `assinador/target/assinador.jar`.

### 1.2. (Opcional) Instalar o JAR manualmente

Use este passo apenas se quiser forçar o uso da versão compilada localmente em vez da versão baixada automaticamente.

**Windows (PowerShell):**
```powershell
New-Item -ItemType Directory -Force "$HOME\.hubsaude"
Copy-Item "assinador\target\assinador.jar" "$HOME\.hubsaude\assinador.jar"
```

**Linux / macOS:**
```bash
mkdir -p ~/.hubsaude
cp assinador/target/assinador.jar ~/.hubsaude/assinador.jar
```

O CLI procura o JAR em `~/.hubsaude/assinador.jar`. A versão local tem precedência sobre o download automático.

---

## 2. Usando o CLI `assinatura`

Todos os comandos são executados na **raiz do projeto** com `go run ./cmd/assinatura`.

> **Por que `go run`?** O binário `assinatura` não está instalado no sistema — `go run` compila e executa na hora, sem precisar instalar nada. É a forma recomendada para desenvolvimento.
>
> **Quer usar sem `go run`?** Compile uma vez com `go build -o assinatura.exe ./cmd/assinatura` e use `.\assinatura.exe` no lugar de `go run ./cmd/assinatura` em todos os exemplos abaixo.

### 2.1. Ver a versão

```powershell
go run ./cmd/assinatura version
go run ./cmd/assinatura --version
```

Saída esperada em desenvolvimento:
```
assinatura dev windows/amd64
```

Em binários publicados via GitHub Releases, `dev` é substituído pela versão rastreável da tag e SHA curto do commit (ex.: `v0.2.0+a3f1c9b`).

---

### 2.2. Criar uma assinatura digital

```bash
go run ./cmd/assinatura sign --local --content "documento.pdf"
```

Parâmetros:

| Parâmetro | Obrigatório | Descrição |
|---|---|---|
| `--content` | Sim | Conteúdo a ser assinado (nome do arquivo, texto, etc.) |
| `--token` | Não | PIN ou token de autenticação do dispositivo criptográfico |
| `--local` | Não | Invoca o JAR diretamente em vez de usar o modo servidor HTTP |

Saída esperada:
```
✔ Assinatura criada com sucesso
  Assinatura : MOCKED_SIGNATURE_BASE64_==
  Mensagem   : Assinatura criada com sucesso
```

O valor em **Assinatura** é o que você usa no próximo passo para validar.

---

### 2.3. Validar uma assinatura digital

```powershell
go run ./cmd/assinatura validate --local --content "documento.pdf" --signature "MOCKED_SIGNATURE_BASE64_=="
```

Parâmetros:

| Parâmetro | Obrigatório | Descrição |
|---|---|---|
| `--content` | Sim | O mesmo conteúdo usado ao criar a assinatura |
| `--signature` | Sim | A assinatura retornada pelo comando `sign` |
| `--local` | Não | Invoca o JAR diretamente em vez de usar o modo servidor HTTP |

Saída quando a assinatura é **válida**:
```
✔ Assinatura VÁLIDA
  Mensagem : Assinatura é válida
```

Saída quando a assinatura é **inválida**:
```
✘ Assinatura INVÁLIDA
  Mensagem : Assinatura é inválida
```

---

## 3. Fluxo completo de exemplo

```bash
# Passo 1: criar assinatura
go run ./cmd/assinatura sign --local --content "contrato.pdf"

# Saída:
# ✔ Assinatura criada com sucesso
#   Assinatura : MOCKED_SIGNATURE_BASE64_==

# Passo 2: validar a assinatura obtida
go run ./cmd/assinatura validate --local --content "contrato.pdf" --signature "MOCKED_SIGNATURE_BASE64_=="

# Saída:
# ✔ Assinatura VÁLIDA

# Passo 3: testar com assinatura errada
go run ./cmd/assinatura validate --local --content "contrato.pdf" --signature "qualquer-coisa-errada"

# Saída:
# ✘ Assinatura INVÁLIDA
```

---

## 4. Modos de invocação

O CLI suporta dois modos de chamar o `assinador.jar`:

| Modo | Como ativar | Quando usar |
|---|---|---|
| **Servidor HTTP** (padrão) | nenhuma flag extra | Múltiplas operações em sequência — elimina o overhead de inicialização da JVM a cada chamada |
| **Local (direto)** | `--local` | Operações esporádicas ou quando o servidor não está rodando |

> O servidor HTTP será iniciado com `assinatura start` (Sprint 3). Por enquanto use sempre `--local`.

---

## 5. Ajuda integrada

Cada comando tem ajuda embutida com exemplos de uso:

```powershell
go run ./cmd/assinatura --help
go run ./cmd/assinatura sign --help
go run ./cmd/assinatura validate --help
go run ./cmd/simulador --help
go run ./cmd/simulador start --help
```

O `--help` mostra a descrição do comando, os parâmetros disponíveis e exemplos práticos de uso.

---

## 6. Decisões de arquitetura

As decisões técnicas relevantes do projeto estão documentadas como ADRs em [`docs/adr/`](adr/):

| ADR | Decisão |
|---|---|
| [001](adr/001-modo-servidor-como-padrao.md) | Modo servidor HTTP como padrão de invocação |
| [002](adr/002-go-para-cli.md) | Go como linguagem dos CLIs |
| [003](adr/003-porta-padrao.md) | Porta padrão 8080 para o assinador.jar |

---

## 7. Solução de problemas

### "não foi possível conectar ao servidor do assinador"

O CLI tentou usar o modo HTTP (padrão) mas o servidor não está rodando. Use `--local` para invocar o JAR diretamente:

```powershell
go run ./cmd/assinatura sign --local --content "documento.pdf"
```

### "assinador.jar não encontrado e não foi possível baixá-lo"

O CLI tentou baixar o JAR automaticamente mas não conseguiu. Verifique a conexão com a internet e tente novamente. Como alternativa, compile localmente seguindo o [Passo 1.1](#11-opcional-compilar-o-assinadorjar-localmente).

### "JDK não encontrado e download automático falhou"

O CLI tentou baixar o JRE 21 automaticamente mas não conseguiu. Instale manualmente o [Eclipse Temurin 21](https://adoptium.net/) e adicione ao PATH. Na próxima execução, o CLI usará o Java do sistema.

### Caracteres estranhos na saída (Windows)

Execute o comando abaixo antes de usar o CLI para corrigir a codificação do terminal:

```powershell
chcp 65001
```
