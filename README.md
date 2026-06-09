# Sistema Runner

> Trabalho Prático — Implementação e Integração de Software (UFG, 2026-1)

O **Sistema Runner** é um conjunto de ferramentas CLI que permite executar aplicações Java sem que o usuário precise configurar o ambiente Java. É composto por três aplicações:

| Componente      | Tecnologia | Descrição                                                                   |
|-----------------|------------|-----------------------------------------------------------------------------|
| `assinatura`    | Go 1.25    | CLI para criação e validação de assinaturas digitais simuladas              |
| `simulador`     | Go 1.25    | CLI para gerenciar o ciclo de vida do Simulador do HubSaúde                |
| `assinador.jar` | Java 21    | Serviço de assinatura/validação (modo local via `java -jar` e modo HTTP)   |

---

## Estrutura do Repositório

```
.
├── cmd/
│   ├── assinatura/                  ← binário do CLI de assinatura
│   │   ├── main.go
│   │   └── cmd/
│   │       ├── root.go
│   │       ├── version.go
│   │       ├── version_test.go
│   │       ├── sign.go
│   │       └── validate.go
│   └── simulador/                   ← binário do CLI do simulador
│       ├── main.go
│       └── cmd/
│           ├── root.go
│           ├── start.go             ← Sprint 4 (US-03.1)
│           ├── stop.go              ← Sprint 4 (US-03.2)
│           └── status.go            ← Sprint 4 (US-03.2)
│
├── internal/
│   ├── invoker/
│   │   └── local.go                 ← invocação do assinador.jar via java -jar
│   ├── jdk/
│   │   └── jdk.go                   ← detecção e provisionamento do JDK 21
│   └── release/
│       └── .gitkeep                 ← download de artefatos (Sprint 4)
│
├── assinador/                       ← projeto Java/Maven
│   └── src/
│       ├── main/java/com/kyriosdata/assinador/
│       │   ├── Main.java            ← ponto de entrada do jar
│       │   ├── service/
│       │   │   ├── SignatureService.java      ← interface principal
│       │   │   └── FakeSignatureService.java  ← implementação simulada
│       │   └── domain/
│       │       ├── SignRequest.java
│       │       ├── SignatureResponse.java
│       │       └── ValidateRequest.java
│       └── test/java/com/kyriosdata/assinador/
│           └── service/
│               └── FakeSignatureServiceTest.java
│
├── docs/
│   └── manual-usuario.md            ← entregável Sprint 4
│
├── orientacoes_trabalho_disciplina/ ← material fornecido pelo professor
│   ├── especificacao.md
│   ├── design.md
│   ├── diagramas/
│   └── docs/
│
├── .github/
│   └── workflows/
│       ├── build.yml                ← CI: testa em push para main
│       └── release.yml              ← CD: publica release ao criar tag v*
│
├── go.mod
├── go.sum
└── .gitignore
```

---

## Pré-requisitos

| Ferramenta | Obrigatório | Descrição |
|---|---|---|
| [Go 1.25+](https://go.dev/dl/) | Sim | Para compilar e rodar os CLIs |
| Java 21 JDK | Não | Baixado automaticamente pelo CLI se ausente |
| [Maven 3.9+](https://maven.apache.org/) | Apenas dev | Para compilar o `assinador.jar` localmente |
| [Cosign](https://docs.sigstore.dev/cosign/system_config/installation/) | Não | Para verificar artefatos do GitHub Releases |

---

## Como compilar

### CLIs Go

```bash
# Compilar ambos os binários
go build ./cmd/assinatura
go build ./cmd/simulador

# Verificar
go vet ./...

# Testar
go test ./...
```

### assinador.jar (Java)

```bash
cd assinador/
mvn clean package
# Gera: assinador/target/assinador.jar
```

### Build multiplataforma

```bash
VERSION=v0.1.0

# assinatura
GOOS=linux   GOARCH=amd64 go build -ldflags "-X github.com/kyriosdata/runner/cmd/assinatura/cmd.version=${VERSION}" -o dist/assinatura-${VERSION}-linux-amd64   ./cmd/assinatura
GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/kyriosdata/runner/cmd/assinatura/cmd.version=${VERSION}" -o dist/assinatura-${VERSION}-windows-amd64.exe ./cmd/assinatura
GOOS=darwin  GOARCH=amd64 go build -ldflags "-X github.com/kyriosdata/runner/cmd/assinatura/cmd.version=${VERSION}" -o dist/assinatura-${VERSION}-darwin-amd64  ./cmd/assinatura

# simulador
GOOS=linux   GOARCH=amd64 go build -ldflags "-X github.com/kyriosdata/runner/cmd/simulador/cmd.version=${VERSION}" -o dist/simulador-${VERSION}-linux-amd64   ./cmd/simulador
GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/kyriosdata/runner/cmd/simulador/cmd.version=${VERSION}" -o dist/simulador-${VERSION}-windows-amd64.exe ./cmd/simulador
GOOS=darwin  GOARCH=amd64 go build -ldflags "-X github.com/kyriosdata/runner/cmd/simulador/cmd.version=${VERSION}" -o dist/simulador-${VERSION}-darwin-amd64  ./cmd/simulador
```

---

## Como rodar (desenvolvimento)

Os CLIs são executados com `go run` a partir da raiz do repositório. Não é necessário instalar o binário.

> O `assinador.jar` e o JDK 21 são baixados automaticamente na primeira execução se não estiverem disponíveis localmente.

### CLI `assinatura`

```powershell
# Ver a versão
go run ./cmd/assinatura version
go run ./cmd/assinatura --version

# Assinar (modo local — invoca java -jar diretamente)
go run ./cmd/assinatura sign --local --content "documento.pdf"
go run ./cmd/assinatura sign --local --content "documento.pdf" --token "meu-pin"

# Validar assinatura (modo local)
go run ./cmd/assinatura validate --local --content "documento.pdf" --signature "MOCKED_SIGNATURE_BASE64_=="

# Assinar via servidor HTTP (padrão — requer assinatura start, Sprint 3)
go run ./cmd/assinatura sign --content "documento.pdf"

# Ajuda
go run ./cmd/assinatura --help
go run ./cmd/assinatura sign --help
```

### CLI `simulador`

```powershell
# Iniciar o Simulador do HubSaúde (Sprint 4)
go run ./cmd/simulador start

# Parar o simulador (Sprint 4)
go run ./cmd/simulador stop

# Verificar status (Sprint 4)
go run ./cmd/simulador status
```

### Compilar o binário (opcional)

Se preferir rodar sem `go run`:

```powershell
go build -o assinatura.exe ./cmd/assinatura
.\assinatura.exe sign --local --content "documento.pdf"
```

---

## Roadmap — Sprints

### ✅ Sprint 1 — Fundação e CI/CD

- [x] Módulo Go inicializado (`github.com/kyriosdata/runner`)
- [x] CLI `assinatura` com Cobra — comando `version` funcionando
- [x] CLI `simulador` com Cobra — estrutura `start/stop/status` criada
- [x] GitHub Actions: build multiplataforma a cada push
- [x] GitHub Actions: release com SemVer ao criar tag `v*`
- [x] Checksums SHA-256 por binário
- [x] Assinatura com Cosign/Sigstore (keyless OIDC)

### 🔄 Sprint 2 — Assinatura Simulada (modo local)

- [x] Projeto Java (`assinador/`) com estrutura Maven
- [x] Interface `SignatureService` e implementação `FakeSignatureService`
- [x] Testes unitários Java passando
- [ ] **US-02.2** — Validação completa dos parâmetros no `assinador.jar` (formato + presença)
- [ ] **US-02.3** — Validação de parâmetros e simulação do fluxo de `validate`
- [ ] **US-01.2** — Comandos `sign` e `validate` no CLI com `--help` documentado
- [ ] **US-01.3** — Invocação real do `assinador.jar` via `java -jar` com testes de integração
- [ ] **US-01.4** — Exibição legível dos resultados no terminal
- [ ] **US-04.1** — Download automático do JDK 21 quando ausente

### ⏳ Sprint 3 — Modo Servidor HTTP e PKCS#11

- [ ] **US-02.4** — Endpoints `POST /sign` e `POST /validate` no `assinador.jar`
- [ ] **US-02.5** — Integração com dispositivo criptográfico via PKCS#11 (SoftHSM2)
- [ ] **US-01.5** — `assinatura start` inicia o jar como servidor, registra PID em `~/.hubsaude/`
- [ ] **US-01.6** — CLI usa modo HTTP por padrão quando servidor está ativo
- [ ] **US-01.7** — Detecção de instância já em execução (health check)
- [ ] **US-01.8** — `assinatura stop` encerra o servidor
- [ ] **US-01.9** — `--timeout <min>` para encerramento automático por inatividade

### ⏳ Sprint 4 — CLI Simulador e Entrega Final

- [ ] **US-03.1** — `simulador start` verifica portas e inicia o `simulador.jar`
- [ ] **US-03.2** — `simulador stop` e `simulador status` com registro em `~/.hubsaude/`
- [ ] **US-03.3** — Pipeline CI/CD gerando binários do `simulador` multiplataforma
- [ ] **US-03.4** — Download automático do `simulador.jar` do GitHub Releases com checksum
- [ ] Manual do usuário (`docs/manual-usuario.md`) completo

---

## Entregáveis obrigatórios

| # | Entregável | Status |
|---|-----------|--------|
| 1 | Código-fonte CLI `assinatura` (Go, multiplataforma) | 🔄 Em andamento |
| 2 | Código-fonte `assinador.jar` (Java 21, Maven) | 🔄 Em andamento |
| 3 | Testes (unitários, integração, aceitação, cenários de erro) | 🔄 Em andamento |
| 4 | Documentação (manual de usuário, guia técnico, exemplos) | ⏳ Sprint 4 |
| 5 | Especificação com diagramas C4 | ✅ Fornecida pelo professor |
| 6 | Binários pré-compilados Win/Linux/macOS via GitHub Releases | ✅ Automatizado |
| 7 | Código-fonte CLI `simulador` (Go, multiplataforma) | ⏳ Sprint 4 |

---

## Dados locais (`~/.hubsaude/`)

```
~/.hubsaude/
├── jdk/            ← JDK 21 provisionado automaticamente (Sprint 2)
├── simulador.jar   ← baixado dinamicamente (Sprint 4)
└── state.json      ← PID, porta e estado dos processos (Sprint 3)
```

---

## CI/CD

| Workflow | Gatilho | O que faz |
|----------|---------|-----------|
| `build.yml` | Push na `main` | `go vet`, `go test`, cross-compile para 3 plataformas |
| `release.yml` | Tag `v*` | Build com versão injetada, checksums SHA-256, assinatura Cosign, publica no GitHub Releases |

### Verificar autenticidade de um artefato

```bash
cosign verify-blob \
  --certificate assinatura-v1.0.0-linux-amd64.pem \
  --signature   assinatura-v1.0.0-linux-amd64.sig \
  --certificate-identity "https://github.com/<usuario>/runner/.github/workflows/release.yml@refs/tags/v1.0.0" \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  assinatura-v1.0.0-linux-amd64
```

---

## Arquitetura (C4 — Contêineres)

```
Usuário
  ├──(CLI)──► assinatura (Go 1.24)
  │               ├──(java -jar)──► assinador.jar (Java 21) ──(PKCS#11)──► Dispositivo Criptográfico
  │               └──(HTTP)──────► assinador.jar (modo servidor)
  └──(CLI)──► simulador (Go 1.24)
                  └──(HTTP)──────► Simulador do HubSaúde (sistema externo)
```

Diagramas detalhados: [`orientacoes_trabalho_disciplina/diagramas/`](orientacoes_trabalho_disciplina/diagramas/)

---

## Referências

- [Especificação do projeto](orientacoes_trabalho_disciplina/especificacao.md)
- [Design e arquitetura C4](orientacoes_trabalho_disciplina/design.md)
- [Plano de implementação (Sprints)](orientacoes_trabalho_disciplina/docs/plano-revisitado-v2.md)
- [Casos de uso — Criar Assinatura (FHIR HubSaúde)](https://fhir.saude.go.gov.br/r4/seguranca/caso-de-uso-criar-assinatura.html)
- [Casos de uso — Validar Assinatura (FHIR HubSaúde)](https://fhir.saude.go.gov.br/r4/seguranca/caso-de-uso-validar-assinatura.html)
