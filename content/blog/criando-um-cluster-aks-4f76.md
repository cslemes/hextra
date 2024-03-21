---
title: Cluster Kubernetes na Azure com Terraform
excerpt: Criando um Cluster Kubernetes na Azure com Terraform e AKS   Neste tutorial, vamos explorar...
author: Cristiano Lemes
coverImage: https://media.dev.to/cdn-cgi/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2F66evzntvh18s6gmysklp.jpg
---
Neste tutorial, vamos explorar o processo passo a passo para criar um cluster Kubernetes na plataforma Azure usando o Terraform e o Azure Kubernetes Service (AKS). Este guia é especialmente útil para iniciantes que desejam iniciar sua jornada com o Kubernetes na nuvem.
 

***Conteúdo*** 
-  [Criação da conta Azure](#1)
-  [Configurando o Ambiente Local](#2)
-  [Criando o Cluster](#3)
-  [Destruindo o Cluster](#4)
 
##  1. Criação da conta Azure <a  id=1>  </a> 
  
###  Criando uma conta na Azure.

Se você ainda não possui uma conta Azure, pode aproveitar os benefícios de recursos gratuitos por 12 meses, além de um crédito de $200 para uso durante 30 dias. Embora os recursos deste tutorial não estejam incluídos na camada gratuita de 12 meses, os $200 podem ser suficientes seguindo os passos e removendo os recursos após o tutorial.
 
1. Acesse a página inicial do [Azure](https://azure.microsoft.com/pt-br/).
2. Clique em "Conta Gratuita" no canto superior direito.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/8tz4ctzt98mda5oapy67.png"  alt="Azure home"  width="60%"  height="auto">
3. Selecione "Experimente Gratuitamente" e siga as instruções para criar uma nova conta ou fazer login com uma conta existente da Microsoft.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/o5q7rnusauglhrcdx2pq.png"  alt="Azure crie conta"  width="60%"  height="auto">
4. Na tela de login, você pode logar com uma conta microsoft já existe, ou criar uma conta nova, se você criar uma conta nova, pode usar um email atual de qualquer provedor.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/66yydv47xvepst3xfw6m.png"  alt="Azure crie Senha"  width="40%"  height="auto">
5. Preencha a região e data de nascimento.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/h3x5t4k6s3j40gsqu3z0.png"  alt="Azure crie conta"  width="40%"  height="auto">
6. Proceda com a verificação do email.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/2pmz18kfdb8o974qoh3c.png"  alt="Azure crie conta"  width="40%"  height="auto">
7. Preencha os detalhes do perfil, e confirme o numero de telefone, marque a caixa verificar identidade por telefone e enviar sms, ou receber ligação.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/z7lu6nngjtwlyu7ut0a0.png"  alt="Azure perfil"  width="40%"  height="auto">
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ly1b1q95w9pqoheetfsd.png"  alt="Azure crie conta"  width="40%"  height="auto">
8. Preencha os dados do endereço, confirme os termos de uso e clique avançar.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/64xntklhtbn1iqy0pabb.png"  alt="Azure crie conta"  width="40%"  height="auto">
9. Adicione os dados do cartão de crédito, e clique em increver-se.
<img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/icpvlrqyyq0slmz0fmtr.png"  alt="Azure crie conta"  width="40%"  height="auto">

Caso você não tenha crédito liberado, ou já tenha passado dos limites ou do tempo do uso dos crédito você terá que habilitar o pagamento pelo uso, Pay as you go, para poder acessar os recursos.

##  2. Configurando o ambiente local <a  id=2>  </a>

Neste tutorial, usaremos o Windows 11, mas as etapas podem ser seguidas em sistemas Linux ou macOS.  

###  1. Instalando o Azure Cli

Existem várias maneiras de instalar o Azure CLI. Para instalação local:
 Windows
Instalar usando o Msi
[Azure Cli Msi](https://aka.ms/installazurecliwindowsx64)
 
Instalar usando um gerenciador de pacotes,
Usando Scoop

```powershell
scoop install azure-cli
```

Usando Winget
```powershell
winget install -e --id Microsoft.AzureCLI

```
Ubuntu
```Bash
curl  -sL  https://aka.ms/InstallAzureCLIDeb | sudo  bash
```
Mac OS
```bash
brew  update && brew  install  azure-cli
```

###  2. Instalando o Kubectl

Windows

Baixando o binario diretamente, você colocar em qualquer diretorio quie preferir, depois adicione esse diretorio a variavel de sistema PATH , para poder executar de qualquer pasta na linha de comando.

[Kubectl Windows Binario]https://dl.k8s.io/release/v1.29.2/bin/windows/amd64/kubectl.exe

Ulizando um gerenciador de pacotes.
Scoop
```powershell
scoop install kubectl
```
Winget
```powershell
winget install -e --id Kubernetes.kubectl
```
Linux
Faça o donwload do binario

```bash
curl  -LO  "https://dl.k8s.io/release/$(curl  -L  -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
```
Configure as permissões e mova para pasta bin.
```bash
sudo  chmod  +x  kubectl
sudo  mv  ./kubectl  /usr/local/bin/
```
Mac OS
Instalando usando homebrew
```bash
brew  install  kubectl
```

### 3. Instalando o Terraform

*Windows*
Baixando o binario, descompacte o zip e coloque em uma pasta de sua preferencia, lembrando que deve por a pasta na variavel PATH.

[Terrafom zip](https://releases.hashicorp.com/terraform/1.7.4/terraform_1.7.4_windows_amd64.zip)

Usando gerenciador de pacotes.
*Scoop*
```powershell
$ scoop  install  terraform
```
*Winget*
```powershell
$ winget install --id=Hashicorp.Terraform -e
```
*Linux*
Utilizando o gerenciador de pacotes (Ubuntu)

```bash
$ wget -O- https://apt.releases.hashicorp.com/gpg | sudo  gpg  --dearmor  -o  /usr/share/keyrings/hashicorp-archive-keyring.gpg
$ echo  "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release  -cs) main" | sudo  tee  /etc/apt/sources.list.d/hashicorp.list
$ sudo  apt  update && sudo  apt  install  terraform
```
*Mac OS*
```bash
$ brew  tap  hashicorp/tap
$ brew  install  hashicorp/tap/terraform
```

4. Instalando GIT

*Windows*

Utilizando o instalador.
Execute o instalador, pode manter todas as configurações padrões.
[Git Windows Instaler](https://github.com/git-for-windows/git/releases/download/v2.44.0.windows.1/Git-2.44.0-64-bit.exe)
Utilizando gerenciador de pacotes
*Scoop*
```powershell
$ scoop install git
```
*Winget*
```powershell
$ winget install --id Git.Git -e --source winget
```
*Linux*
A maioria das distros já vem com git, utilize o gerenciador de pacotes para instalar caso necessário.
*Ubuntu*
```bash
$ apt-get  install  git
```
*Mac OS*
```bash
$ brew  install  git
```
5. Agora precisamos configurar o ambiente local para acessar o Azure.

Fazendo login na conta Azure usando o azure cli, na linha de comando digite.
```powershell
$ az login
```
Ele vai abrir o navegador padrão e solicitar as credenciais do azure, entre com as credencias ele automaticamente vai configurar a linha de comando para poder acessar, os dados ficaram salvo na pasta ~./azure.

###  3. Criando o Cluster <a  id=3>  </a>

Para criar o cluster, siga estes passos:

1. Vamos baixar o template AKS da HashiCorp para servir como ponto de partida, clone o repositorio usando o git.
```powershell
$ git clone https://github.com/hashicorp/learn-terraform-provision-aks-cluster aks-cluster
```
Para ajustar os arquivos você pode utilizar o editor de código de sua preferencia, eu estarei utilizando o [Visual Studio Code](https://code.visualstudio.com/).

Com o Vscode aberto, vá em arquivo abrir pasta, escolha a pasta que acabou criar no passo anterior, ou navegue até a pasta usando a linha de comando e digite code .
 <img  src="https://dev-to-uploads.s3.amazonaws.com/uploads/articles/cx4gk0xmy4wt5841sjyx.png"  alt="Azure crie conta"  width="auto"  height="auto">
 
2. Gerando credenciais para o Terraform

Vá na linha de comando
```powershell
$ az ad sp create-for-rbac --skip-assignment
```
Ele vai gerar os dados que vamos usar para configurar o acesso ao Azure utilizando Terraform
Exemplo da saida.

```powershell
{

"appId": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
"displayName": "azure-cli-2019-04-11-00-46-05",
"password": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
"tenant": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}

```
3. Editando os arquivos do terraform
Edite o arquivo terraform.tfvar e adicione os valores de appId e password que você obteve da saida do comando anterior.
Edit o arquivo aks-cluster.tf em kubernetes altere para uma versão do kubernetes que seja suportado pelo aks, no momento "1.28.3", para verificar as versões do kubernetes disponiveis use o comando abaixo, eu vou manter a zona com West US2.

```powershell
$ az aks get-versions --location westus2 --output table

Kubernetes Version Upgrades
------------------- -----------------------

1.28.3 None available
1.28.0  1.28.3
1.27.7  1.28.0, 1.28.3
1.27.3  1.27.7, 1.28.0, 1.28.3
1.26.10  1.27.3, 1.27.7
1.26.6  1.26.10, 1.27.3, 1.27.7
```

Cada nó do cluster AKS precisa ter no minimo 2vCPU e 4GB de memoria, neste caso vou utilizar a maquina mais barata disponivel com essas configurações que é a Standard_B2s.
para verficar os preços de maquinas pode verificar no site do [azure](https://azure.microsoft.com/en-us/pricing/details/virtual-machines/windows/).
O arquivo aks-cluster.tf final fica assim:

```terraform
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "random_pet" "prefix" {}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "default" {
  name     = "${random_pet.prefix.id}-rg"
  location = "West US 2"

  tags = {
    environment = "Demo"
  }
}

resource "azurerm_kubernetes_cluster" "default" {
  name                = "${random_pet.prefix.id}-aks"
  location            = azurerm_resource_group.default.location
  resource_group_name = azurerm_resource_group.default.name
  dns_prefix          = "${random_pet.prefix.id}-k8s"
  kubernetes_version  = "1.28.3"

  default_node_pool {
    name            = "default"
    node_count      = 2
    vm_size         = "Standard_B2s"
    os_disk_size_gb = 30
  }

  service_principal {
    client_id     = var.appId
    client_secret = var.password
  }

  role_based_access_control_enabled = true

  tags = {
    environment = "Demo"
  }
}

```
4. Agora vamos inicializar o terrform
na linha de comando digite:

```powershell
 
$ terraform init

Initializing the backend...
Initializing provider plugins...

- Reusing previous version of hashicorp/random from the dependency lock file
- Reusing previous version of hashicorp/azurerm from the dependency lock file
- Installing hashicorp/random v3.5.1...
- Installed hashicorp/random v3.5.1 (signed by HashiCorp)
- Installing hashicorp/azurerm v3.67.0...
- Installed hashicorp/azurerm v3.67.0 (signed by HashiCorp)
  

Terraform has made some changes to the provider dependency selections recorded
in the .terraform.lock.hcl file. Review those changes and commit them to your
version control system if they represent changes you intended to make.
Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```
Com isso o terraform instala todos os modulos necessários.

5.  Agora vamos rodar o terraform plan, o terraform plan vai simular a criação do recursos na cloud, é recomendado salvar o plano em um arquivo para posterior execução do apply.

```powershell
$ terraform plan --out plan1

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the

following symbols:

+ create
Terraform will perform the following actions:
......
Plan: 3 to add, 0 to change, 0 to destroy.
Changes to Outputs:

+ kubernetes_cluster_name = (known after apply)
+ resource_group_name = (known after apply)
```

6. Com o plano sem erros podemos rodar o apply.
```powershell
$ terraform apply "plan1"
---

Apply complete! Resources: 3 added, 0 changed, 0 destroyed.
Outputs:

kubernetes_cluster_name = "apt-tortoise-aks"
resource_group_name = "apt-tortoise-rg"
```

7. Capture os dados do output para os dados do cluster, e execute o comando abaixo para configurar o kubectl.

```powershell
az aks get-credentials --resource-group apt-tortoise-rg --name apt-tortoise-aks --file ~/.kube/config
```
8. Agora o cluster está rodando, vamos verifcar com o kubectl.
```
$ kubectl get nodes

NAME STATUS ROLES AGE VERSION
aks-default-27467048-vmss000000 Ready agent 3m17s v1.28.3
aks-default-27467048-vmss000001 Ready agent 3m26s v1.28.3
```

9. Vamos fazer um pequeno deploy para testar o nosso cluster de um servidor Nginx.

```powershell

$ kubectl create deployment website --replicas 3 --image nginx
$ kubectl expose deployment/website --type="LoadBalancer" --port 80

```

10. O Loadbalancer do Azure vai disponibilizar um IP publico para o serviço, e nesta configuração padrão não tem regras de segurança, então ele vai estar acessivel de qualquer lugar.
Verificando os pods, criado com a quantidade de replicas que solicitamos.

```powershell

$ kubectl get svc

NAME TYPE CLUSTER-IP EXTERNAL-IP PORT(S) AGE
kubernetes ClusterIP 10.0.0.1 <none> 443/TCP 13m
website LoadBalancer 10.0.195.30 x.x.x.x 80:30707/TCP 45s

$ kubectl get pods
NAME READY STATUS RESTARTS AGE
website-6784674d46-6xddr 1/1 Running 0 8m42s
website-6784674d46-7nsl6 1/1 Running 0 8m42s
website-6784674d46-vn5r4 1/1 Running 0 8m42s

```
11. Efetuando um curl no ip fornecido você deve receber a pagina padrão do nginx

```powershell
$ curl 4.246.49.237
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>
<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>
<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

##  4. Destruindo o Cluster

<a  id=4>  </a>
 
1. Agora para não ficar gerando custos desnecessários para um ambiente de testes, vamos destruir o cluster, sempre que quiser é só criar o cluster novamente rodando o apply.


```powershell
$ terrafom destroy

Plan: 0 to add, 0 to change, 3 to destroy.

Changes to Outputs:
- kubernetes_cluster_name = "apt-tortoise-aks" -> null
- resource_group_name = "apt-tortoise-rg" -> null

Do you really want to destroy all resources?
Terraform will destroy all your managed infrastructure, as shown above.
There is no undo. Only 'yes' will be accepted to confirm.
Enter a value: yes
```

### Conclusão
Este guia fornece uma visão geral abrangente para criar e gerenciar um cluster Kubernetes na Azure. Lembre-se de revisar cuidadosamente e adaptar as etapas conforme necessário para atender às suas necessidades específicas e às últimas práticas recomendadas.


### Referências

- [Instalação Azure Cli](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli)
- [Instalação Kubectl](https://kubernetes.io/pt-br/docs/tasks/tools/)
- [Instalação Terraform](https://developer.hashicorp.com/terraform/install)
- [Git](git-scm.com)

Ilustração
- [3D ico](https://3dicons.co/)
- [Simple icons](https://simpleicons.org/)
