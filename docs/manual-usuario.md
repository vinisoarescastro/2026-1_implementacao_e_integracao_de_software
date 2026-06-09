# Manual do Usuário — Sistema Runner

O Sistema Runner é um conjunto de ferramentas de linha de comando (CLI) que permite criar e validar assinaturas digitais sem precisar conhecer os detalhes do ambiente Java.

## Pré-requisitos

| Ferramenta | Versão mínima | Como verificar |
|---|---|---|
| Go | 1.25 | `go version` |
| Java JDK | 21 | `java -version` |
| Maven | 3.9 | `mvn -version` |

---

## 1. Configuração inicial (primeira vez)

Estes passos só precisam ser feitos uma vez.

### 1.1. Compilar o assinador.jar

Na raiz do projeto, execute:

```bash
mvn -f assinador/pom.xml package -q
```

O comando gera o arquivo `assinador/target/assinador.jar`.

### 1.2. Instalar o JAR no diretório do Runner

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

O CLI sempre procura o JAR em `~/.hubsaude/assinador.jar`.

---

## 2. Usando o CLI `assinatura`

Todos os comandos abaixo devem ser executados na raiz do projeto com `go run ./cmd/assinatura`.

> **Dica:** após compilar o binário com `go build -o assinatura ./cmd/assinatura`, você pode trocar `go run ./cmd/assinatura` por `./assinatura` em todos os exemplos.

### 2.1. Ver a versão

```powershell
go run ./cmd/assinatura version
go run ./cmd/assinatura --version
```

Saída esperada:
```
assinatura dev windows/amd64
```

Ambas as formas funcionam. Em binários publicados via GitHub Releases, `dev` é substituído pela versão da tag (ex.: `v0.2.0`).

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

Cada comando tem ajuda embutida:

```bash
go run ./cmd/assinatura --help
go run ./cmd/assinatura sign --help
go run ./cmd/assinatura validate --help
```

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

### "assinador.jar não encontrado"

O JAR não está em `~/.hubsaude/`. Repita o [Passo 1.2](#12-instalar-o-jar-no-diretório-do-runner).

### "não foi possível localizar ou provisionar o JDK"

O Java 21 não está instalado ou não está no PATH. Instale o JDK 21 (recomendado: [Eclipse Temurin](https://adoptium.net/)) e tente novamente.

### Caracteres estranhos na saída (Windows)

Execute o comando abaixo antes de usar o CLI para corrigir a codificação do terminal:

```powershell
chcp 65001
```
