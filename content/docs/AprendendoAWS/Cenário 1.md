---
title: Aprendendo AWS
date: 2024-02-15T23:09:29.755Z
type: docs
---

Nesse artigo vamos aprender os conceitos básicos da AWS na prática utilizando um cenário real de migração para cloud.

- Vamos Utilizar: IAM, EC2, EBS, S3, RDS, S3, ELB, AWS Backup, VPC,  Application Migration Service
#### Avaliando os Requisitos

*Qual o Négocio*?
  - Empresa de varejo, com lojas físicas, tipo americanas e casas Bahia.
  - 10 Filiais, sendo 5 em Shopping
  - Funciona de 08:00 as 22:00 de segunda a sábado
  - Lojas de Shopping abrem de 14:00 as 22:00 no Domingo
  - 5 pontos de vendas por loja (50 total)
  - 20 Usuários no Escritório
   
Para deixar definir quais sistema a empresa vai usar, resolvi seguir as informações da [EAESP, da FGV](https://eaesp.fgv.br/producao-intelectual/pesquisa-anual-uso-ti), que elabora uma pesquisa anualmente para avaliar os recursos TI em uso nas empresas brasileiras.

![Sistemas Operacionais Servidor](images/eaesp-so.png)
*O SO mais usado em servidores é Windows.*

![Evolução Tendencia](images/eaesp-trend.png)
*O Linux tem crescido, mas tá longe de ser o padrão*

![Sistema Operacional Micro](images/eaesp-so-ws.png)
*Nos desktops o Windows mantem a dominância, ainda vai chegar o ano do Linux do desktop*

![Sistemas ERP](images/eaesp-erp.png)
*Sistemas ERP mais usados*

![Banco de Dados](images/easp-database.png)
*Banco de dados*

Levando em consideração esses dados, nossa empresa vai seguir esse padrão, maioria dos servidores Windows, Desktop Windows, 

Para o ERP Vamos seguir a média dos requisitos dos ERP mais usados, lembrando que estamos tomando como referencia somente, não estamos simulando a instalação e migração real dessa aplicação.

Nossa empresa tem apenas 20 usuários para o ERP, portanto:
- Cada usuário consome 200MB então 20 x 200 4000MB + 8196MB(Para o SO) = 12GB no servidor de aplicação.
- 8vCPU, ele recomenda 4 CPU físicas para até 100 usuários, como vamos usar virtualização, então 8vCPU, 4 cores + HyperThread. 

| Servidor de Aplicação | Virtual        |
| --------------------- | -------------- |
| Sistema Operacional   | Windows Server |
| CPU                   | 8vCPU          |
| Mémoria RAM           | 12 GB          |
| Disco SO + PageFile   | 50GB           |
| Disco Dados           | 200 GB         |
O cliente (é uma aplicação cliente servidor), no momento é instalada em cada desktop da empresa.
Aplicação Cliente

- CPU: 5% a 30% de uso do processador.
- Memória RAM: 100 MB a 1 GB de memória.
- Rede: 1 MB a 10 MB de largura de banda por hora de uso.

Para o banco de dados, vamos usar o MySQL.

| Servidor de Aplicação | Bare Metal |
| --------------------- | ---------- |
| Sistema Operacional   | Debian     |
| CPU                   | 4 CPU      |
| Mémoria RAM           | 32 GB      |
| Disco SO              | 100GB      |
| Disco Dados           | 500 GB     |

***O que tem na empresa***
 - 1 Servidor de Aplicação Windows (Todas roles da aplicação estão nesse servidor)
 - 1 Banco de dados Linux
 - 1 Servidor de Arquivos Windows
 - 1 Servidor de Impressão Windows
 - 1 Active directory Windows
 - 1 Zabbix Linux
 - 1 DHCP 
 - 1 DNS Windows

***Diagrama da Empresa***


```mermaid


flowchart TB
classDef blue fill:#40e0d0
classDef green fill:#50c878
classDef yellow fill:#ffd700
classDef darkblue fill:#5c9dca
classDef olive fill:#556b2f
classDef orange fill:#e47025
classDef brow fill:#634217
classDef red fill:#ce5a57
classDef darkgrey fill:#444c5c
classDef darkyellow fill:#e1b16a

subgraph Datacenter
    direction TB
    ST[(Storage)]:::darkyellow
	FS<-.->ST
    subgraph Rede Privada      
      direction TB
      RT{{"Roteador"}}:::green   
      SW1[Switch TOR]:::orange
      FW[Firewall]:::red   
      RT -.-> SW1
      FW -.-> RT     
    end  
    subgraph HyperVisor
    direction TB      
       DNSP(DNS Privado)     
       RDP[\"Terminal Server"\]:::olive
       FS[\"Servidor de Arquivos"\]:::olive
       APP[\"Aplicação Web"\]:::olive
	   SW1 --> DNSP & RDP & FS & APP         
	end
    subgraph Bare Metal
      DB[("Banco de Dados")]:::darkgrey
      DB1[("Banco de Stand BY")]:::darkgrey
      APP-.->DB
      RDP-.->DB
      DB<-.->DB1
	end
end

subgraph Escritorio
   direction TB
   RTE{{"Roteador"}}:::green 
   SWE["Switch"]:::orange
   USRE>"Usuarios"]:::darkblue
   USRE-.->SWE-.->RTE
end
subgraph Filiais
    direction TB
	RTF{{"Roteador"}}:::green 
	SWF["Switch"]:::orange
	USRF>"Usuarios"]:::darkblue
	USRF-.->SWF-.->RTF
end


MPLS("MPLS")
IE(("Internet"))
MPLS-.->FW 
RTF-.->MPLS
RTE-.->MPLS
IE<-.->FW
 

```

#### *Porque migrar?*

Uma parte importante de migrar ou não para cloud, requer analises de objetivos do negócio, analise de custos, arquitetura de aplicativos, segurança, desempenho.

*Quais gaps precisam ser resolvidos*
- Negócio: A empresa está crescendo e o ambiente atual não suporta as demandas, umas das soluções é adquirir mais hardware
- Segurança física: A única proteção de acesso aos servidores é uma chave.
- Energia: A sala de dados possui somente em no-break que suporta 2 horas sem energia, e possui somente um fornecedor de energia externa.
- Disponibilidade
- Escalabilidade
- Backup são feitos localmente

#### O que é cloud?




#### *Datacenter Tiers*

Os "tiers" de data centers referem-se a um sistema de classificação que descreve a disponibilidade e confiabilidade das instalações de um data center. Existem quatro níveis principais, cada um representando um padrão de infraestrutura e redundância, conforme definido pelo [Uptime Institute](https://uptimeinstitute.com/) :

- Tier 1: O Tier 1 é o nível básico e oferece a menor disponibilidade. Geralmente, possui uma infraestrutura simples, sem redundância de componentes críticos. Pode ser propenso a interrupções para manutenção e atualizações.
- Tier 2: O Tier 2 oferece um pouco mais de confiabilidade que o Tier 1, com alguma redundância em componentes críticos, como fontes de energia e sistemas de refrigeração. No entanto, ainda pode ser interrompido para manutenção planejada.
- Tier 3: O Tier 3 é projetado para fornecer uma maior disponibilidade do que os níveis anteriores. Possui redundância em todos os componentes críticos e permite a manutenção sem interrupção dos serviços. Geralmente, possui N+1 de redundância, o que significa que há backup completo para todos os componentes essenciais.
- Tier 4: O Tier 4 é o nível mais alto de disponibilidade e confiabilidade. Ele oferece redundância completa em todos os aspectos, incluindo energia, refrigeração, conectividade de rede e segurança. Além disso, é projetado para suportar falhas individuais sem interrupção dos serviços. Os data centers Tier 4 são os mais caros de construir e operar, mas oferecem o mais alto nível de garantia de tempo de atividade.

Lista de Datacenters certificados pelo [Instituto Uptime](https://uptimeinstitute.com/uptime-institute-awards/country/id/BR), a maioria são Tier 3.

![Datacenters Brasil](images/datacentersbr.png)
*Mapa de datacenters certificados espalhados pelo Brasil*

Verificando na lista, vamos ver que temos empresas especializadas em datacenter, bancos e empresas governamentais em sua maioria, até grandes empresas não desfrutam de ter um próprio datancer, inclusive cloud providers tem seus servidores em colocation com em empresas como Equinix, Tivit e Ascenty por exemplo.
Sendo assim a nossa empresa tem duas opções disponíveis para alcançar suas metas dentro do seu budget, colocation e cloud.

 
#### ***Regiões e Zonas  AWS***
  
As regiões e zonas da Amazon Web Services (AWS) são parte da infraestrutura global da AWS para hospedar serviços em nuvem. 

1\. Regiões: As regiões da AWS são áreas geográficas separadas que consistem em várias zonas de disponibilidade. Cada região é composta por dois ou mais data centers que são isolados fisicamente e estão localizados em áreas distintas para aumentar a resiliência e a disponibilidade dos serviços. As regiões da AWS estão localizadas em todo o mundo e permitem que os clientes implantem aplicativos em locais geograficamente diversos para melhorar a resiliência e a latência.
    
2\. Zonas de Disponibilidade: As zonas de disponibilidade são data centers isolados dentro de uma região que são conectados por redes de baixa latência e alta largura de banda. Cada zona de disponibilidade é projetada para ser independente das outras zonas, com infraestrutura de energia, refrigeração e rede próprias. Isso significa que uma falha em uma zona de disponibilidade não afetará as outras.
    
Em resumo, as regiões da AWS representam áreas geográficas distintas, enquanto as zonas de disponibilidade são data centers isolados dentro dessas regiões.  

![Aws Global](images/aws-global.png)
[*Infraestrutura Global da  Aws*](https://aws.amazon.com/pt/about-aws/global-infrastructure/regions_az/)



### Criando a Infraestrutura na AWS

#### ***1\. Criando usuários usando o IAM***

Vamos pular a parte de criar a conta na AWS, para não ficar tão longo, já que o processo é simples, qualquer dúvida, pode olhar na [documentação oficial](https://docs.aws.amazon.com/accounts/latest/reference/welcome-first-time-user.html)
Vamos criar um usuário no Aws Console, não é recomendado usar a conta root para tarefas guarde ela para caso haja algum problemas nas contas Admin e precise recupera-las.

##### 1.1. Regras básicas sobre usuários, grupos e politicas.

- As permissões na AWS são definidas por politicas
- Há politicas pré definidas com diferentes tipos de acesso nos serviços da AWS, como leitura, escrita e acesso total.
- Uma política pode ser aplicada a um Grupo ou diretamente a um usuário(para facilitar a gestão melhor aplicar sempre a grupos) 
- Um grupo pode receber N politicas
- Um usuário pode participar de N grupos
- Um grupo não pode ser membro de outro grupo
- Quando um usuário pertence a mais de um grupo, as politicas aplicadas aos grupos que ele pertence são somadas.

![IAM](images/iam.png)
*Relação entre grupos, políticas e usuários
##### 1.2. Criando um Grupo e definindo permissões

- Vá até Página inicial do console
- Na barra de pequisa digite **IAM**
- Em services, **IAM**, marque a estrela se quiser deixar nos favoritos.
- Click em **User Groups**
- Click em **Create Group**
- Escolha um nome para o grupo, esse nome é exclusivo somente na sua organização.
- Marque a política padrão, **Administrator Access**. 
	*Há várias predefinições de políticas, e você também pode criar novas, o recomendado é ser mais específico possível e habilitar somente o recurso que cada time precisa para efetuar o seu trabalho, no caso vamos escolher administradores, para seguir esse tutorial.*
- Click em **Create group**

![Create group](videos/Create-Group.gif)
##### 1.3. Criando o usuário e adicionando ao grupo criado

- Estando em IAM, click em **Users**
- Click em **Create User**
- Em ***User details**, escreva o nome de usuário, ele é exclusivo somente dentro da organização. Não é recomendado o uso de usuários genéricos, então crie um usuário para cada pessoa do time.
-  Marque **Provide user access to the AWS Management Console** 
- Escolha **I want to create an IAM user.**
- Click em **Next**
- Em **User groups** escolha o grupo que acabamos de criar, no caso AWS-Admins
- Click em **Next**
- Click em **Create User**
- Conta de Usuário criado com êxito, na tela de criação podemos obter a senha do usuário, que deixamos em criar automaticamente, e o login direto para a console.
- Efetue logoff com da conta root e logue com o usuário IAM criado.

![Create User](videos/Create-User.mp4)
#### ***2\. Configurando o Budget na AWS***

Importante saber que na AWS não conseguimos travar os gastos com os serviços, uma maneira de controlar os gastos é criando budgets, com os budgets podemos definir um valor e receber alertas quando ele for atingido.

##### 2.1. Criando um Orçamento custo zero.

- Vá até a Página inicial do Aws console
- Escreva **billing** na barra de pesquisa
- Click em **Billing and Cost Management**
- Marque a estrela para deixar nos favoritos (opcional)
- Click em **Budgets**
- Click em **Create a Budget**
- Escolha *Use a template (simplified)*
- Escolha **My Zero-Spend Budget**
- Digite o email para onde serão envidas as notificações.
- Click em **Create budget**
- Seu orçamento **My Zero-Spend Budget** foi criado.

![Zero Budget](videos/ZeroSpendBudget.mp4)

##### 2.2. Criando um orçamento mensal estipulando um valor  

- Estando em  **Billing and Cost Management**
- Click em **Budgets**
- Click em **Create Budget**
- Escolha *Use a template (simplified)*
- Escolha **Monthly cost budget**
- Vou colocar 5 doláres
- Digite o email para onde serão enviadas as notificações
- Click em **Create budget**
- Seu orçamento My Monthly Cost Budget foi criado.

![Cost Budget](videos/MonthlyCostBudget.mp4)

##### 2.3. Analisando Custos

Em explorador de custos você tem um relatório com os gastos, com opção de vários filtros, como intervalo de datas e nome de serviços entre outros, se você usou algum serviço Free Tier, você pode visualizar o quanto usou dele em nível gratuito.

![Report Budget](images/budgets-reports.png)


#### ***3\. Criando a infraestrutura de redes***

O primeiro item de infraestrutura que vamos criar é a rede, você pode criar outros itens não tendo a rede criada, mas nesse caso a AWS vai criar automaticamente uma rede padrão para este serviço.
Vamos criar duas subnets uma vai ter ip público, vai poder ser  e outra vai ter apenas ips privados. 
Nossa infra não será muito grande, por isso vamos uma VPC /22 vai ser suficiente.
-- Dúvidas sobre CIDRs, pode olhar nessa documentação da [Digital Ocean](https://www.digitalocean.com/community/tutorials/understanding-ip-addresses-subnets-and-cidr-notation-for-networking).

Os 3 Primeiros ips de cada subrede são reservados pela AWS, o primeiro ip é para o VPC Router (Gateway), o segundo é o DNS da Amazon, o terceiro está reservado para uso futuro.
Lembrando que o primeiro ip da Subnet é o id de rede, é o ultimo é o ip de Broadcast, que também não podem serem usados.

3.1 Overview serviços da VPC

- ***VPC:*** Significa Virtual Private Cloud, permite que você crie uma rede virtual, permitindo isolamento entre recursos.
- ***Subnet:*** É uma subdivisão do intervalo de ip da VPC, serve para organizar, criar rotas, regras de segurança, permitindo uma gerencia do tráfego na VPC.
- ***RouteTable:*** É uma tabela de roteamento, que por padrão já adiciona automaticamente todas rotas da VPC, você vai adicionar manualmente rota para fora da vpc, como outras vpcs, serviços da Aws como s3 e rds, e rotas para internet.
- ***Nat Gateway:*** É um roteador que não está voltado para internet, não fornecendo ip publico para as instancias, ele é totalmente gerenciado pela AWS, possuindo alta disponibilidade, escalabilidade e segurança nativos.
- ***Internet Gateway:*** É um roteador que fornecesse acesso a internet, tanto atribuindo um ip público diretamente a instância, como fornecendo internet atraves de NAT pelo Nat gateway.
- ***Security Groups:*** Age como firewall, permitindo que você controle o tráfego de redes com base de regras, é atribuído a nivel de instância,  configurando regras baseado em endereços IP de origem ou destino, portas e protocolos, como TCP, HTTP. Uma importante saber, é que as regras são statefull, significa que ele salva o estado das conexões, portanto quando você crie uma regra, as repostas de saídas já são liberadas automaticamente.
- ***NACL:*** Siginifica Network Access Control List, ele adiciona mais uma camada de segurança, ele é associado a subrede, quer dizer, as regras são aplicadas a todas as instancias dentro da mesma subrede por padrão já vem com tudo liberado, do mesmo modo que security group você pode criar regras com ips, portas e protocolos de origem, mas você pode tanto permitir quando negar um acesso, e seguem a ordem que você estabelece numericamente, parando de processar quando tem o match. Importante saber que ele é stateless, não salva o estado das conexões, quer dizer, quando criar uma regra de entrada tem que criar também a regra de saída. 
- ***VPC Endpoints:*** VPC Endpoints fornecem conexões diretas a serviços gerenciados pela AWS, como S3, DynamoDB e SNS, sem a necessidade de roteamento pela Internet pública. Eles garantem uma comunicação segura e eficiente entre a VPC e esses serviços, contribuindo para uma arquitetura mais robusta e segura na nuvem.


Diagrama da rede.

```mermaid
flowchart TB;
subgraph Region["sa-east-1 São Paulo"]
	direction TB
    subgraph VPC["VPC 10.1.0.0/22"]
	    direction TB
        Subrede1 --> Router{{Roteador}}
        Router --> IEGateway1[Internet Gateway]    
        EDP["VPC Endpoint"]
	    
	    
		subgraph AZA["Zona de Disponibilidade A"]
		  direction TB
		  
		  NG1[NAT Gateway]
	      subgraph Subrede1["PublicSubnetA"]
	      direction TB
		      RT1{{"PublicRouteTable"}} 
		      Inst1["Rede 10.1.0.0/24"]
	          Inst2["Gateway 10.1.0.1"]
	          
	          RT1 ~~~ Inst1 ~~~ Inst2 ~~~ NG1
          end
          subgraph Subrede2["PrivateSubnetA"]
          direction TB
	          RT2{{"PrivateRouteTable"}}
	          Inst3["Rede 10.1.1.0/24"]
	          Inst4["Gateway 10.1.1.1"]
	          RT2 ~~~ Inst3 ~~~ Inst4
	       end
        end
    end    
    S3[("S3")]
end      
IEGateway1[Internet Gateway] --> IE((Internet))
```



##### 3.1. Criando uma VPC

- Vá até a Página inicial do Aws console
- Escreva **vpc** na barra de pesquisa
- Click em **VPC**
- Marque a estrela para deixar nos favoritos (opcional)
- Click em **Create VPC**
- Marque **VPC only**
- Em Name tag, coloque o nome da rede vpc-01
- Deixe marcado **IPv4 CIDR block**
- Em **IPv4 CIDR4** coloque **10.1.0.0/22**
- Deixe marcado **No IPv6 CIDR block**
- Click em **Create VPC**

![Create VPC](videos/CreateVPC.mp4)
##### 3.2. Criando Subnets

Vamos criar duas subnets na AZ a, uma pública e outra privada.

- Estando em VPC
- Click em **Subnets**
- 
- Em **VPC ID** escolha a VPC que acabamos de criar
- Em **Subnet settings** click em Add New Subnet, 2 vezes.
- Primeira subnet
	- Subnet name: PublicSubnetA
	- Availability Zone: sa-east-1a
	- IPv4 subnet CIDR block: 10.1.0.0/24
- Segunda subnet
	- Subnet name: PrivateSubnetA
	- Availability Zone: sa-east-1a
	- IPv4 subnet CIDR block: 10.1.1.0/24

<video scr=..videos/CreateSubnetAZa.mp4/>

##### 3.3. Configurando a subnet pública

***3.3.1.  Habilitando Subnet para receber ips públicos por padrão ***
- Estando em VPC
- Click em **Subnets**
- Na lista de subnets , click com botão direito em PublicSubnetA, escolha **edit subnet settings**
- Em Auto-assign IP settings, marque a opção **Enable auto-assign public IPv4 address**
- Click em save



****3.3.2. Criando um internet gateway***
- Estando em VPC
- Click em **Internet gateways**
- Click no botão **Create internet gateways**
- Em Name tag, vamos colocar **IGW01**
- Click em **Create internet gateway** 

***3.3.3. Anexando o internet gateway a subnet**
- Estando em VPC
- Em **Internet gateways**
- Na lista de Internet gateways , click com botão direito em IGW01, escolha **attach to VPC**


#### *4\. Criando uma VM usando EC2*


#### ***5\. Conectando o ambiente on premises com a cloud***

Nessa etapa vamos criar uma VPN Ipsec entre a Digital Ocean e a AWS para simular a conexão entre o ambiente on premises com a cloud.

#### ***6\. Migrando servidores para AWS***

AWS Application migration service.

#### 7\. Migrando servidor de arquivos para S3.

#### 8\. Migrando Banco de Dados para Amazon RDS

#### 9.\ Implementando Backups

#### 10.\ Implementando Monitoramento CloudWatch

#### *11.\ Redirecionando os serviços para aws usando Route 53*



