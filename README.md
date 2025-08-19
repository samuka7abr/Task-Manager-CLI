# TaskManSm7 - Gerenciador de Tarefas em CLI

Um gerenciador de tarefas simples e eficiente para linha de comando, desenvolvido em Go. Permite gerenciar suas tarefas diárias com prioridades, datas de vencimento e visualização em calendário.

## 🛠️ Tecnologias Utilizadas

![Go](https://skillicons.dev/icons?i=go)
![Linux](https://skillicons.dev/icons?i=linux)

## 🚀 Instalação

### Pré-requisitos

- Go 1.24.5 ou superior
- Git

### Instalação via Make

```bash
# Clone o repositório
git clone https://github.com/samuka7abr/Task-Manager-CLI.git
cd Task-Manager-CLI

# Instale o projeto
make install
```

## 📖 Como Usar

### Comandos Básicos

```bash
# Ver ajuda geral
taskmansm7 --help

# Adicionar uma tarefa
taskmansm7 add "Reunião com cliente"

# Adicionar tarefa com prioridade
taskmansm7 add "Relatório mensal" alta

# Adicionar tarefa com data específica
taskmansm7 add "Entrega projeto" media 15 12 2024

# Listar tarefas do dia
taskmansm7 day

# Listar tarefas de uma data específica
taskmansm7 day 20 12 2024

# Listar tarefas da semana
taskmansm7 week

# Visualizar calendário do mês
taskmansm7 cal

# Editar uma tarefa (modo interativo)
taskmansm7 edit

# Deletar uma tarefa por ID
taskmansm7 del 1
```

## 🎨 Funcionalidades

### Sistema de Prioridades
- **Alta** (vermelho): Tarefas urgentes e importantes
- **Média** (amarelo): Tarefas normais
- **Baixa** (verde): Tarefas de baixa urgência

### Visualizações

#### Lista de Tarefas
Exibe tarefas em formato tabular com:
- ID da tarefa
- Prioridade (colorida)
- Data de vencimento
- Nome da tarefa

#### Calendário com Mapa de Calor
- Visualização mensal
- Densidade de tarefas por dia (mapa de calor)
- Dia atual destacado
- Legendas explicativas

### Persistência de Dados
- As tarefas são salvas automaticamente em `~/.taskman_cli.json`
- Dados persistem entre sessões
- Formato JSON legível

## 🛠️ Desenvolvimento

### Estrutura do Projeto
```
Task-Manager-CLI/
├── main.go          # Código principal
├── Makefile         # Comandos de build e instalação
├── go.mod           # Dependências Go
├── go.sum           # Checksums das dependências
└── README.md        # Este arquivo
```

### Dependências
- `github.com/spf13/cobra`: Framework para CLI
- `github.com/spf13/viper`: Gerenciamento de configuração

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 👨‍💻 Autor

**Samuel Abrão**
- GitHub: [@samuka7abr](https://github.com/samuka7abr)
- Portfólio: [portifolio-lyart-three-23.vercel.app](https://portifolio-lyart-three-23.vercel.app/)

