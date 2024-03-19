
Nesse artigo vamos aprender os conceitos básicos da AWS na prática utilizando um cenário real de migração para cloud.

### 1- Migrando uma SMB de on premises para AWS

- Vamos Utilizar: IAM, EC2, EBS, S3, RDS, S3, ELB, AWS Backup, VPC,  Application Migration Service
#### Avaliando os Requisitos

*Qual o Négocio*?
  - Empresa de varejo, com lojas físicas, tipo americanas e casas Bahia.
  - 10 Filiais, sendo 5 em Shopping
  - Funciona de 08:00 as 22:00 de segunda a sábado
  - Lojas de Shopping abrem de 14:00 as 22:00 no Domingo
  - 5 pontos de vendas por loja (50 total)
  - 20 Usuários no Escritório
   

Para definir quais sistema a empresa vai usar, resolvi seguir as informações da [EAESP, da FGV](https://eaesp.fgv.br/producao-intelectual/pesquisa-anual-uso-ti). 


![[Pasted image 20240310160830.png]]
*O SO mais usado em servidores é Windows.*

![[Pasted image 20240310161050.png]]
*O Linux tem crescido, mas tá longe de ser o padrão*

![[Pasted image 20240310161353.png]]
*Nos desktops o Windows mantem a dominância, ainda vai chegar o ano do Linux do desktop*


![[Pasted image 20240310161703.png]]
*Sistemas ERP mais usados*



![[Pasted image 20240310171637.png]]
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


O que é cloud?




#### *Datacenter Tiers*

Os "tiers" de data centers referem-se a um sistema de classificação que descreve a disponibilidade e confiabilidade das instalações de um data center. Existem quatro níveis principais, cada um representando um padrão de infraestrutura e redundância, conforme definido pelo [Uptime Institute](https://uptimeinstitute.com/) :

- Tier 1: O Tier 1 é o nível básico e oferece a menor disponibilidade. Geralmente, possui uma infraestrutura simples, sem redundância de componentes críticos. Pode ser propenso a interrupções para manutenção e atualizações.
- Tier 2: O Tier 2 oferece um pouco mais de confiabilidade que o Tier 1, com alguma redundância em componentes críticos, como fontes de energia e sistemas de refrigeração. No entanto, ainda pode ser interrompido para manutenção planejada.
- Tier 3: O Tier 3 é projetado para fornecer uma maior disponibilidade do que os níveis anteriores. Possui redundância em todos os componentes críticos e permite a manutenção sem interrupção dos serviços. Geralmente, possui N+1 de redundância, o que significa que há backup completo para todos os componentes essenciais.
- Tier 4: O Tier 4 é o nível mais alto de disponibilidade e confiabilidade. Ele oferece redundância completa em todos os aspectos, incluindo energia, refrigeração, conectividade de rede e segurança. Além disso, é projetado para suportar falhas individuais sem interrupção dos serviços. Os data centers Tier 4 são os mais caros de construir e operar, mas oferecem o mais alto nível de garantia de tempo de atividade.

Lista de Datacenters certificados pelo [Instituto Uptime](https://uptimeinstitute.com/uptime-institute-awards/country/id/BR), a maioria são Tier 3.

![[Pasted image 20240308202347.png ]]
*Mapa de datacenters certificados espalhados pelo Brasil*

Verificando na lista, vamos ver que temos empresas especializadas em datacenter, bancos e empresas governamentais em sua maioria, até grandes empresas não desfrutam de ter um próprio datancer, inclusive cloud providers tem seus servidores em colocation com em empresas como Equinix, Tivit e Ascenty por exemplo.
Sendo assim a nossa empresa tem duas opções disponíveis para alcançar suas metas dentro do seu budget, colocation e cloud.

 
#### ***Regiões e Zonas  AWS***
  
As regiões e zonas da Amazon Web Services (AWS) são parte da infraestrutura global da AWS para hospedar serviços em nuvem. 

1\. Regiões: As regiões da AWS são áreas geográficas separadas que consistem em várias zonas de disponibilidade. Cada região é composta por dois ou mais data centers que são isolados fisicamente e estão localizados em áreas distintas para aumentar a resiliência e a disponibilidade dos serviços. As regiões da AWS estão localizadas em todo o mundo e permitem que os clientes implantem aplicativos em locais geograficamente diversos para melhorar a resiliência e a latência.
    
2\. Zonas de Disponibilidade: As zonas de disponibilidade são data centers isolados dentro de uma região que são conectados por redes de baixa latência e alta largura de banda. Cada zona de disponibilidade é projetada para ser independente das outras zonas, com infraestrutura de energia, refrigeração e rede próprias. Isso significa que uma falha em uma zona de disponibilidade não afetará as outras.
    
Em resumo, as regiões da AWS representam áreas geográficas distintas, enquanto as zonas de disponibilidade são data centers isolados dentro dessas regiões.  

![[Pasted image 20240308211253.png]]
[*Infraestrutura Global da  Aws*](https://aws.amazon.com/pt/about-aws/global-infrastructure/regions_az/)



### Criando a Infraestrutura na AWS

#### ***1\. Criando usuários usando o IAM***

Vamos pular a parte de criar a conta na AWS, para não ficar tão longo, já que o processo é simples, qualquer dúvida, pode olhar na [documentação oficial](https://docs.aws.amazon.com/accounts/latest/reference/welcome-first-time-user.html)
Vamos criar um usuário no Aws Console, não é recomendado usar a conta root para tarefas guarde ela para caso haja algum problemas nas contas Admin e precise recupera-las, logado com sua conta root.

Regras básiocas sobre usuários, grupos e politicas

- As permissões na AWS são definidas por politicas
- Há politicas pré definidas com diferentes tipos de acesso nos serviços da AWS, como leitura, escrita e acesso total.
- Uma política pode ser aplicada a um Grupo ou diretamente a um usuário(para facilitar a gestão melhor aplicar sempre a grupos) 
- Um grupo pode receber N politicas
- Um usuário pode participar de N grupos
- Um grupo não pode ser membro de outro grupo
- Quando um usuário pertence a mais de um grupo, as politicas aplicadas aos grupos que ele pertence são somadas.

![[Pasted image 20240309122254.png]]

Seguindo na prática

***[1\. Vá até Página inicial do console | Página inicial do console | sa-east-1](https://sa-east-1.console.aws.amazon.com/console/home?region=sa-east-1)***

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/0f35a42b-8fef-4b5d-8725-452927605faf/1/50/50?1)

2\.  Na barra de pequisa digite IAM

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/8c6b6f49-3d11-4dfb-a616-985d9167cf1d/1.7318232316998/37.351585995834/0.3940110323089?1)

3\.  Selecione em services, IAM, marque a estrela se quiser deixar nos favoritos.


![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/91ca81e6-b838-4b68-82c6-2decc2a33dfd/2.5/42.087021549191/12.687155240347?1)

4\. Primeiro vamos criar um Grupo e definir as permissões.

![](https://dubble-prod-01.s3.amazonaws.com/assets/214a52a3-e776-42cf-8857-d730d56e65d8.png?1)

5\. Click em Grupos de usuários

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/8305768c-504e-4603-9ef9-405fa408bc8f/2.5/2.8318584070796/25.978462333469?1)

6\. Click em Criar grupo
![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/5de9bc9b-2b4b-4a02-ac80-c2e2cf1d0ef0/2.5/89.321537524198/11.899133175729?1)

7\. Escolha um nome para o grupo, esse nome é exclusivo somente na sua organização.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/d74f5034-e455-4991-97ec-ca83a9915e56/1.76/52.212389380531/24.481220410695?1)


8\. Marque a política padrão, Administrator Access. Há várias predefinições de politica, o recomendado é ser mais específico possível e habilitar somente o recurso que cada Time precisa para efetuar o seu trabalho, no caso vamos escolher administradores, para seguir esse tutorial. 

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/a9147d8b-1c0a-4659-875a-d5fd0add2128/2.5/30.088495575221/38.744419780277?1)

9\. Click on Criar grupo
![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/b395edf9-d00c-4243-870f-15b507a6d025/2.5/89.616520974488/95.50828069472?1)

10\. Com o grupo criado agora vamos criar o usuário, click em Usuários

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/b2a8fe4f-6529-4d6f-a5db-2885c20020b0/2.5/2.8318584070796/28.342528527322?1)

11\. Click em Criar usuário

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/e122f807-3e1e-4515-b285-4add0651280b/2.5/89.321537524198/13.475177304965?1)

12\. Escreva o nome de usuário, ele é exclusivo somente dentro da organização. Não é recomendado o uso de usuários genéricos, então crie um usuário para cada pessoa do time.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/5b5dfb96-3846-474d-bc63-9a13036ea9b2/1.2384615384615/62.389380530973/24.166011584848?1)

13\. Marque Fornecer acesso ao Console

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/664b45dd-3b58-4476-ae8c-b288404a9174/2.5/33.628318584071/29.839770450096?1)

14\. Escolha Quero criar um usuário do IAM

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/741ee4d6-c672-40fd-8837-c1b2f0b3f5ce/1.5372595406584/64.778763762618/46.992385659845?1)

15\. Click em Próximo

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/be75a482-2a60-4fb4-9429-00fa78e8429f/2.5/91.091449028623/93.459423326713?1)

16\. Em grupo de usuários escolha o grupo que acabamos de criar

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/cd23ae92-1916-4692-979f-918ab24f7d6f/2.5/33.628318584071/64.827950119127?1)

17\. Click em Próximo

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/c9d4a10d-eeef-4620-bc6a-8df106f6aef5/2.5/91.091449028623/77.304966192343?1)

18\. Click em Criar usuário

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/daf22ac9-c905-4e2a-b437-ea535d3f6887/2.5/91.091449028623/69.345943339703?1)

19\. Conta Usuário criado com êxito, na tela de criação podemos obter a senha do usuário, que deixamos em criar automaticamente, e o login direto para a console.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/42e34c7e-05d3-43f9-b05d-a65771ddea93/1/50/51.562913634629?1)

20\. Efetue logoff com a conta Root e logue com o usuario IAM criado.

#### ***2\. Configurando o Budget na AWS***

Importante saber que na AWS não conseguimos travar os gastos com os serviços, uma maneira de controlar os gastos é criando budgets, com os budgets podemos definir um valor e receber alertas quando ele for atingido.


1\.Vá até a Página inicial do console

![](https://dubble-prod-01.s3.amazonaws.com/assets/61c9ad5f-b98e-492d-83e2-f9115dde5f99.png?1)

2\. Escreva billing na barra de pesquisa

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/fc114d15-8d80-4d39-a993-489bb21da3e6/1.8666666666667/34.344027523703/0.3940110323089?1)

3\. (opcional) Marque a estrela para deixar nos favoritos
![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/39c31d24-5d6a-4fe1-9bf2-60f17b6790d2/2.5/58.112451019048/19.516680201176?1)

4\. Click em Billing and Cost Management

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/75d0fdb1-7b6d-4834-b7e6-90f28f2faa06/2.5/46.606433372915/19.332808786908?1)

5\. Click em Orçamentos

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/449ff687-7a04-40b9-98a3-45e18a6b795c/2.5/2.5039123630673/69.66115216555?1)

6\. Click em Criar um orçamento

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/eef087a1-c167-4521-a643-6167e06faf77/2.5/79.443792781919/30.627792514714?1)

7\. Escolha Orçamento de gasto zero, e digite o email para onde serão envidas as notificações.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/83592f9b-98f1-415c-8f3f-3acf2d84e2ed/1.4363636363636/52.425665101721/63.567116619378?1)

8\. Click em Criar orçamento

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/5d7add34-c1e6-4db0-bad8-ee163f9e210e/2.5/78.195098494886/92.513796849172?1)

9\.Seu orçamento My Zero-Spend Budget foi criado.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/98315a46-2761-494b-af47-87fae4f83683/1/57.824726134585/51.628580536778?1)

10\. Vamos criar agora um orçamento mensal estipulando um valor,  Click em Criar orçamento.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/82d7b9eb-6bac-4fe7-b16a-d09fb03eb562/2.5/87.428275333697/20.724980299448?1)

11\. Click em Use as configurações recomendadas. Você pode alterar algumas opções de configuração após a criação do orçamento.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/e4f43a3a-76fd-463d-a3c0-938124938dc9/2.5/67.605636190734/27.370633846009?1)

12\. Click em Crie um orçamento mensal que notifique se você exceder ou estiver previsto para exceder o valor do orçamento.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/2f611730-df79-46eb-bb66-88cb641e3853/2.5/76.943143507311/48.647231995543?1)

13\. Vou colocar 5 doláres

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/30407477-5529-492f-b5a2-07dd2f73cb81/1.4363636363636/52.425665101721/58.996586840955?1)

14\. Digite o email para onde serão enviadas as notificações

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/8939fcce-f147-49dd-9664-5e886ce2f490/1.4363636363636/52.425665101721/69.50354790293?1)

15\. Click em Criar orçamento

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/369b3f22-5602-4e82-810a-3b7e4397eb32/2.5/78.195098494886/92.513796849172?1)

16\. Seu orçamento My Monthly Cost Budget foi criado.

![](https://d3q7ie80jbiqey.cloudfront.net/media/image/zoom/4974e0b4-d251-4360-a0d1-b372996507f5/1/57.824726134585/51.628580536778?1)

17.\ Em explorador de custos você tem um relatório com os gastos, com opção de vários filtros, como intervalo de datas e nome de serviços entre outros, se você usou algum serviço Free Tier, você pode visualizar o quanto usou dele em nível gratuito.

![[Pasted image 20240309130128.png]]


#### ***3\. Criando a infra de redes VPC***

O primeiro item de infraestrutura que vamos criar é a rede, você pode criar outros itens não tendo a rede criada, mas nesse caso a AWS vai criar automaticamente uma rede padrão para este serviço.
A Vamos criar duas subnets uma vai ter ip público, vai poder ser  e outra vai ter apenas ips privados. 
Nossa infra não será muito grande, por isso vamos uma VPC /22 vai ser suficiente.
Dúvidas sobre CIDRs, pode olhar nessa documentação da [Digital Ocean](https://www.digitalocean.com/community/tutorials/understanding-ip-addresses-subnets-and-cidr-notation-for-networking).

Os 3 Primeiros ips de cada subrede são reservados pela AWS, o primeiro ip é para o VPC Router (Gateway), o segundo é o DNS da Amazon, o terceiro está reservado para uso futuro.
Lembrando que o primeiro ip da Subnet é o id de rede, é o ultimo é o ip de Broadcast, que também não podem serem usados.

Diagrama da rede.

```mermaid
flowchart TB;
    subgraph VPC["VPC 10.1.0.0/22"]
        Subrede1 --> IEGateway1[Internet Gateway]
	    Subrede2 --> NATGateway1[NAT Gateway] 
        Subrede3 --> IEGateway2[Internet Gateway]
        Subrede4 --> NATGateway2[NAT Gateway] 
        subgraph AZA["AZ A"]
	      subgraph Subrede1["PublicSubnetA"]
		      Instancia1["Rede 10.1.0.0/24"]
	          Instancia2["Gateway 10.1.0.1"]
          end
          subgraph Subrede2["PrivateSubnetA"]
	          Instancia3["Rede 10.1.1.0/24"]
	          Instancia4["Gateway 10.1.1.1"]
	       end
        end
		subgraph AZB["AZ B"]
        subgraph Subrede3["PublicSubnetB"]
            Instancia5["Rede 10.1.2.0/24"]
            Instancia6["Gateway 10.1.2.1"]      
        end
        subgraph Subrede4["PrivateSubnetB"]
            Instancia7["Rede 10.1.3.0/24"]
            Instancia8["Gateway 10.1.3.1"]
        end
        end
    end

    

```





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



