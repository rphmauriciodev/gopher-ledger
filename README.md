# 🐹 Gopher Ledger

[![Go Report Card](https://goreportcard.com/badge/github.com/rphmauriciodev/gopher-ledger)](https://goreportcard.com/report/github.com/rphmauriciodev/gopher-ledger)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Gopher Ledger** é um serviço de ledger (livro-razão) de alta performance construído em Go, projetado para processar transações financeiras de forma assíncrona, robusta e escalável.

O projeto utiliza **gRPC** para comunicação de baixa latência e um **Worker Pool** interno para garantir que o processamento de créditos e débitos seja feito de maneira eficiente, respeitando a integridade dos dados via transações ACID no PostgreSQL.

---

## 🚀 Funcionalidades

- **Processamento Assíncrono:** Transações são recebidas via gRPC e enfileiradas para processamento por um pool de workers.
- **Arquitetura Limpa:** Separação clara entre domínio, infraestrutura e handlers.
- **Transações Atômicas:** Garantia de consistência nos saldos das contas.
- **Observabilidade:** Logs estruturados com `zerolog`.
- **Dashboards:** Integração nativa com **Metabase** para visualização de métricas e saldos em tempo real.

---

## 🛠️ Tech Stack

- **Linguagem:** [Go](https://golang.org/) (1.25+)
- **Comunicação:** [gRPC](https://grpc.io/) & [Protocol Buffers](https://developers.google.com/protocol-buffers)
- **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/)
- **Visualização:** [Metabase](https://www.metabase.com/)
- **Containerização:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- **Logs:** [Zerolog](https://github.com/rs/zerolog)
- **Configuração:** [Viper](https://github.com/spf13/viper)

---

## 🏗️ Arquitetura

O Gopher Ledger foi desenhado para suportar alta carga:

1.  **gRPC Handler:** Recebe a requisição `processTransaction`.
2.  **Internal Queue:** A transação é enviada para um canal interno (`chan domain.Transaction`).
3.  **Worker Pool:** N workers (configuráveis) consomem do canal e executam a lógica de negócio (validar saldo, atualizar conta, registrar transação).
4.  **Database Layer:** Persistência garantida via repositórios PostgreSQL com suporte a transações SQL.

---

## 🚦 Como Rodar

### Pré-requisitos

- Docker e Docker Compose instalados.
- Go 1.25+ (se desejar rodar localmente sem Docker).

### Passo a Passo

1.  **Clone o repositório:**
    ```bash
    git clone https://github.com/rphmauriciodev/gopher-ledger.git
    cd gopher-ledger
    ```

2.  **Configure o ambiente:**
    Copie o arquivo de exemplo e ajuste as variáveis se necessário:
    ```bash
    cp .env.example .env
    ```

3.  **Suba os containers:**
    ```bash
    docker-compose up --build
    ```

O servidor gRPC estará rodando na porta definida em `.env` (padrão `:50051`) e o Metabase em `http://localhost:3000`.

---

## 📊 Dashboard (Metabase)

Para visualizar os dados:
1. Acesse `http://localhost:3000`.
2. No setup inicial, crie sua conta de administrador.
3. Quando o Metabase perguntar sobre **"Add your data"**, selecione **PostgreSQL**.
4. Use as seguintes configurações (baseadas no seu `.env`):
   - **Name:** Gopher Ledger
   - **Host:** `db` (importante: use o nome do serviço no Docker, não localhost)
   - **Port:** `5432`
   - **Database name:** `gopher_ledger` (ou o valor de `DB_NAME`)
   - **Username:** seu `DB_USER`
   - **Password:** seu `DB_PASSWORD`
5. Pronto! Agora você pode criar dashboards em cima das tabelas `accounts` e `transactions`.

> [!TIP]
> O Metabase agora utiliza um banco de dados separado (`metabase`) para salvar suas configurações, garantindo que o banco do Ledger (`gopher_ledger`) fique limpo e focado apenas nos dados do negócio.

---

## 📜 API (gRPC)

O contrato do serviço pode ser encontrado em `proto/ledger.proto`:

```protobuf
service LedgerService {
  rpc processTransaction (TransactionRequest) returns (TransactionResponse);
}
```

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

Feito por [Raphael Maurício](https://github.com/rphmauriciodev)
