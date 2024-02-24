# GuicPay Simplificado - Desafio PicPay Backend

Este projeto é uma implementação simplificada do backend do PicPay, desenvolvido como parte do desafio proposto pelo [PicPay](https://github.com/PicPay/picpay-desafio-backend).


## Stack

- **Golang**: Linguagem de programação, compilada, rápida, multi-paradigmas e concorrente.
- **PostgreSQL**: Banco de dados SQL utilizado para armazenar dados persistentes, usufruindo da capacidade de transações atômicas garantindo consistência.
- **Redis**: Sistema de armazenamento em cache atuando como um serviço de lock distribuído.
- **DDD (Domain-Driven Design)**: Metodologia para organizar o código em torno das regras de negócio, onde a modelagem do problema é o mais importante.
- **Clean Architecture**: Estrutura de código que enfatiza a separação de responsabilidades e a independência das camadas e não dependendo de framework.


## Modelagem de Domínio

![Modelagem de Domínio](./assets/model_entity_dark.png)

A imagem acima ilustra a modelagem de domínio do GuicPay simplificado. Cada entidade e sua relação refletem a estrutura fundamental do sistema.


## Arquitetura do Sistema

![Arquitetura do Sistema](./assets/arch_api_dark.png)

A arquitetura do sistema é projetada para ser modular e escalável. Cada camada tem uma responsabilidade específica, facilitando a manutenção e o desenvolvimento contínuo.


## Clean Architecture

![Clean Architecture](./assets/1_O4pMWCi5kZi20SNOR6V33Q.png)

A implementação do GuicPay Simplificado segue os princípios da Clean Architecture. Essa abordagem enfatiza a separação de interesses, facilitando a compreensão do código, a manutenção e a evolução do sistema.


## Como rodar o projeto 🚀

```sh
make docker-run
```

### Health Check

```sh
curl http://localhost:8080/api/ping
```


## Documentação 

Para acessar a documentação OpenAPI basta acessar a rota `/docs` .

