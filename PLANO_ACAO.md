# Plano de Ação — Sistema Runner

Gerado em: 2026-06-08  
Base: auditoria cruzando o repositório com os critérios da disciplina (`orientacoes_trabalho_disciplina/docs/criterios.md` e `especificacao.md`).

---

## Como usar este documento

Cada item tem uma **prioridade** (🔴 Crítico / 🟡 Significativo / 🟢 Polimento) e um **status** ([ ] pendente / [x] feito).  
Resolva na ordem de prioridade. Os itens críticos comprometem a avaliação mesmo que todo o resto esteja certo.

---

## 🔴 Crítico — corrigir antes de qualquer entrega

### C-01 — Inverter a lógica de modo padrão no CLI Go

**Problema:** `Sign()` e `Validate()` em `internal/invoker/local.go` retornam erro quando `local=false`. O modo padrão deveria ser HTTP; `--local` deveria ser a exceção explícita.

**O que fazer:**
- [x] Em `Sign()` e `Validate()`: quando `local=false`, tentar invocar via HTTP (endpoint `POST /sign` e `POST /validate`)
- [x] Implementar `invoker/http.go` com chamada HTTP ao servidor (mesmo que o servidor ainda não exista, o cliente precisa estar correto)
- [x] Quando `--local` for passado, usar `runJar` (comportamento atual)
- [x] Quando o servidor HTTP não responder e `--local` não for passado, exibir erro claro orientando o usuário a usar `--local` ou `assinatura start`

**Referência:** `especificacao.md` US-01 — *"O CLI deve fazer uso do assinador.jar no modo servidor quando não orientado para usar o modo local"*; `criterios.md` E2 — *"Modo servidor é o padrão; modo local deve ser explicitamente ativado"*

---

### C-02 — Adicionar build e testes Java ao CI

**Problema:** `build.yml` executa apenas Go. Se o Java quebrar, o CI não detecta. Viola o princípio de reprodutibilidade.

**O que fazer:**
- [x] Em `build.yml`, adicionar um job `test-java` que execute `mvn -f assinador/pom.xml verify`
- [x] Em `release.yml`, garantir que o JAR é gerado e publicado no GitHub Releases junto com os binários Go
- [x] O job `build` depende de `test-go` e `test-java` via `needs:`

**Referência:** `criterios.md` — *"Reprodutibilidade: qualquer pessoa clona, roda um comando, obtém build e testes verdes"*; `criterios.md` F — *"Distribuição do JAR: artefato único com Main-Class correto"*

---

### C-03 — Implementar download do `assinador.jar` no CLI

**Problema:** `internal/release/` está vazio. O `invoker/local.go` procura o JAR em `~/.hubsaude/` mas nada o coloca lá. Um usuário que clona o repositório não consegue usar `sign --local`.

**O que fazer:**
- [x] Implementar `internal/release/jar.go` com `EnsureJar()`: lê release.json local/remoto, compara versão, baixa só se desatualizado, calcula SHA-256
- [x] Chamar `EnsureJar()` em `invoker/local.go` antes de `runJar`
- [x] Exibir progresso ao usuário durante o download

**Referência:** `especificacao.md` US-03 — estratégia `release.json`; `criterios.md` E1 — *"Erros de execução (JAR ausente) são tratados com mensagens claras"*

---

### C-04 — Implementar download do JDK em `internal/jdk/jdk.go`

**Problema:** `Resolve()` detecta o JDK mas não faz download quando ausente. Retorna erro sem provisionar.

**O que fazer:**
- [x] Implementar download do JRE 21 em `internal/jdk/jdk.go`: detecta plataforma, baixa do Adoptium, extrai `.tar.gz`/`.zip` com proteção contra path traversal
- [x] `Resolve()` chama download automático como terceiro estágio quando JDK não encontrado
- [ ] Teste automatizado que valida `java -version` após provisionamento

**Referência:** `especificacao.md` US-04; `criterios.md` — *"Versões mínimas declaradas e verificadas em runtime com erro amigável"*

---

### C-05 — Criar `release.json` na raiz do repositório

**Problema:** O arquivo não existe no repositório do aluno. O CLI precisa dele para comparar versões e saber de onde baixar o JAR e o JRE.

**O que fazer:**
- [x] Criar `release.json` na raiz com URLs do JAR e do JRE (Eclipse Temurin 21) para as 3 plataformas
- [ ] Automatizar a atualização de `release.json` no `release.yml` a cada nova tag

**Referência:** `especificacao.md` US-03 — *"Busca release.json via URL estável no branch main"*

---

### C-06 — Criar ADRs para decisões não óbvias

**Problema:** Nenhum ADR existe. O `criterios.md` cita explicitamente: *"se não criou nenhuma ADR, por quê?"* — é um critério avaliado.

**O que fazer:**
- [x] Criar pasta `docs/adr/`
- [x] Escrever `docs/adr/001-modo-servidor-como-padrao.md` — por que HTTP é o modo padrão
- [x] Escrever `docs/adr/002-go-para-cli.md` — por que Go foi escolhido
- [x] Escrever `docs/adr/003-porta-padrao.md` — por que porta 8080
- [x] Referenciar as ADRs no `README.md` com tabela linkada

**Template mínimo (1 página cada):**
```
# ADR-NNN — Título

## Contexto
O que motivou esta decisão.

## Decisão
O que foi decidido.

## Consequências
O que fica mais fácil e o que fica mais difícil com essa escolha.
```

**Referência:** `criterios.md` A — *"Decisões registradas (ADRs curtos) onde houve escolha não óbvia"*

---

## 🟡 Significativo — resolver para avaliação completa

### S-01 — Adicionar `.gitattributes`

**O que fazer:**
- [x] Criar `.gitattributes` na raiz com `eol=lf` padrão e `eol=crlf` para `.bat`/`.cmd`

**Referência:** `criterios.md` D — *"Encoding UTF-8 declarado; line endings tratados (.gitattributes)"*

---

### S-02 — Adicionar flag `--version` ao CLI

**Problema:** O CLI tinha subcomando `version`, mas `criterios.md` pede explicitamente a flag `--version`.

**O que fazer:**
- [x] Adicionar `rootCmd.Version = version` em `cmd/assinatura/cmd/root.go`
- [x] Criar `cmd/simulador/cmd/version.go` e adicionar `rootCmd.Version = version` no simulador
- [x] Incluir SHA curto do commit na versão (ex.: `v1.0.0+abc1234`) via ldflags no CI (`release.yml`)

**Referência:** `criterios.md` I — *"Versão acessível via `--version` retornando algo rastreável (tag + SHA curto)"*

---

### S-03 — Adicionar exemplos ao `--help`

**O que fazer:**
- [x] Preencher o campo `Example` de `sign`, `validate`, `simulador start`, `stop` e `status`

**Referência:** `criterios.md` I — *"`--help` que ensina (com exemplos), não que lista flags"*

---

### S-04 — Escrever testes de integração CLI → JAR

**Problema:** Só existe `version_test.go`. Não há testes que invoquem o JAR de verdade.

**O que fazer:**
- [ ] Criar `internal/invoker/local_test.go` com testes que:
  - Verificam se o JAR está disponível (skip se não estiver)
  - Chamam `Sign()` com conteúdo válido e verificam `Result.Signature`
  - Chamam `Validate()` com assinatura correta e verificam `Result.Valid == true`
  - Chamam `Sign()` com conteúdo vazio e verificam erro
- [ ] Marcar com `t.Skip("JAR não disponível")` quando o JAR não estiver presente, para não quebrar CI sem o JAR

**Referência:** `criterios.md` G — *"Testes de contrato CLI → JAR: subprocess real e HTTP real, em ambos os modos"*

---

### S-05 — Compilar e publicar binário do `simulador` no CI/CD

**Problema:** `release.yml` compila apenas `assinatura`. O `simulador` (mesmo que stub) deve estar nos artefatos.

**O que fazer:**
- [x] Em `build.yml`, adicionar cross-compilation do `simulador` nos mesmos alvos
- [x] Em `release.yml`, adicionar build e publicação de `simulador-<versão>-<os>-<arch>`
- [x] Incluir binários do `simulador` no `checksums.txt` e assinar com Cosign

**Referência:** `especificacao.md` seção 7, item 6 — artefatos executáveis incluem `simulador-*`

---

### S-06 — Corrigir `go.mod` para Go 1.25

**O que fazer:**
- [x] Atualizar `go.mod`: `go 1.25`
- [x] Atualizar `build.yml` e `release.yml`: `go-version: '1.25'`

**Referência:** `plano-revisitado-v2.md` — *"CLIs serão desenvolvidos em Go (1.25)"*

---

### S-07 — Completar validação de parâmetros conforme FHIR (US-02.2)

**Problema:** `FakeSignatureService` valida apenas `content`. A US-02.2 exige validação de **todos** os parâmetros conforme a especificação FHIR.

**O que fazer:**
- [ ] Ler os casos de uso referenciados:
  - [Criar Assinatura](https://fhir.saude.go.gov.br/r4/seguranca/caso-de-uso-criar-assinatura.html)
  - [Validar Assinatura](https://fhir.saude.go.gov.br/r4/seguranca/caso-de-uso-validar-assinatura.html)
- [ ] Mapear quais campos são obrigatórios além de `content`
- [ ] Adicionar validações em `FakeSignatureService` para cada campo identificado
- [ ] Adicionar testes unitários para cada nova validação

**Referência:** `especificacao.md` US-02 critério — *"O assinador.jar deve validar todos os parâmetros conforme especificações"*

---

## 🟢 Polimento — para nota máxima

### P-01 — Escrever o manual do usuário

**O que fazer:**
- [x] Substituir o placeholder em `docs/manual-usuario.md` com conteúdo real (pré-requisitos, configuração, exemplos, modos, troubleshooting)
- [ ] Adicionar seção de download dos binários via GitHub Releases com verificação Cosign
- [ ] Adicionar seção do Simulador quando Sprint 4 estiver pronto

**Referência:** `especificacao.md` seção 7, item 4 — *"Manual de usuário para assinatura"*

---

### P-02 — Adotar Conventional Commits daqui em diante

**Problema:** Histórico tem commits com mensagem `"refresh"` sem contexto.

**O que fazer:**
- [ ] Usar prefixos a partir de agora: `feat:`, `fix:`, `chore:`, `test:`, `docs:`, `ci:`
- [ ] Exemplo: `feat: implementar download automático do assinador.jar`
- [ ] Referenciar a issue/história nas mensagens quando aplicável: `feat(US-03.4): baixar simulador.jar do GitHub Releases`

**Referência:** `criterios.md` H — *"Commits atômicos, mensagens no imperativo, idealmente Conventional Commits"*

---

### P-03 — Implementar modo `--verbose` / `--quiet`

**O que fazer:**
- [ ] Adicionar flags globais em `root.go`: `--verbose` e `--quiet`
- [ ] Usar um logger com nível ajustável (ex.: `log/slog` da stdlib do Go 1.21+)
- [ ] Modo verbose: exibir o comando `java -jar` que está sendo executado, headers HTTP, etc.
- [ ] Modo quiet: suprimir tudo exceto o resultado final

**Referência:** `criterios.md` I — *"Logs em nível ajustável; modo `--verbose`/`--quiet` previsível"*

---

### P-04 — Implementar Sprint 3 (modo servidor HTTP)

**O que fazer (após C-01 estar resolvido):**
- [ ] `assinador.jar`: implementar endpoints `POST /sign` e `POST /validate` (US-02.4)
- [ ] CLI: implementar `assinatura start` para iniciar o JAR em modo servidor (US-01.5)
- [ ] CLI: implementar `assinatura stop` para encerrar (US-01.8)
- [ ] CLI: detectar instância em execução via health check HTTP (US-01.7)
- [ ] CLI: auto-shutdown por inatividade com `--timeout` (US-01.9)
- [ ] Registrar PID e porta em `~/.hubsaude/`

---

### P-05 — Implementar Sprint 4 (Simulador HubSaúde)

**O que fazer:**
- [ ] Implementar `simulador start`, `simulador stop`, `simulador status` (US-03.1, US-03.2)
- [ ] Download automático do `simulador.jar` via `release.json` (US-03.4)
- [ ] Verificar disponibilidade da porta 8443 antes de iniciar (US-03.1)
- [ ] Health check via `/api/info` e shutdown via `/shutdown` (US-03.2)

---

## Ordem de execução sugerida

```
Semana atual:
  C-01  Inverter lógica de modo padrão (HTTP default)
  C-02  Adicionar Java ao CI
  C-05  Criar release.json
  C-06  Escrever 3 ADRs
  S-01  Criar .gitattributes
  S-02  Adicionar --version flag
  S-06  Corrigir go.mod para Go 1.25

Próxima semana (Sprint 2 finalização):
  C-03  Download do assinador.jar
  C-04  Download do JDK
  S-03  Exemplos no --help
  S-04  Testes de integração CLI → JAR
  S-05  Compilar simulador no CI/CD
  S-07  Validação FHIR completa

Sprint 3:
  P-04  Servidor HTTP completo (US-01.5 a US-01.9, US-02.4)

Sprint 4:
  P-05  Simulador HubSaúde (US-03.1 a US-03.4)
  P-01  Manual do usuário
  P-03  --verbose / --quiet
```

---

## Rastreabilidade

| Item | User Story | Critério |
|---|---|---|
| C-01 | US-01 | criterios.md E2 |
| C-02 | — | criterios.md (Reprodutibilidade, F) |
| C-03 | US-03, US-04 | criterios.md E1 |
| C-04 | US-04.1 | criterios.md E1 |
| C-05 | US-03 | especificacao.md US-03 |
| C-06 | — | criterios.md A |
| S-01 | — | criterios.md D |
| S-02 | — | criterios.md I |
| S-03 | — | criterios.md I |
| S-04 | US-01, US-02 | criterios.md G |
| S-05 | US-05 | especificacao.md seção 7 |
| S-06 | — | plano-revisitado-v2.md |
| S-07 | US-02.2 | criterios.md E3 |
| P-01 | — | especificacao.md seção 7 item 4 |
| P-02 | — | criterios.md H |
| P-03 | — | criterios.md I |
| P-04 | US-01.5–01.9, US-02.4 | criterios.md E2 |
| P-05 | US-03.1–03.4 | criterios.md E4 |
