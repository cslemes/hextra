---
title: Docker GPU no Windows com WSL2
date: 2024-02-15T23:09:29.755Z
author: Cristiano Lemes
coverImage: ' "https://media.dev.to/cdn-cgi/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Ftavsh0s5mk3gg933pv5l.png"'
---
### Configurando Ambiente de Desenvolvimento com GPU no Windows usando Docker e WSL2

Este guia detalhado ajudará você a configurar um ambiente de desenvolvimento com GPU no Windows usando Docker, WSL2 e as ferramentas relacionadas da Nvidia.

## Passo 1: Verificar os Drivers da GPU

Para começar, verifique se os drivers da sua GPU estão funcionando corretamente. No PowerShell, execute o comando `nvidia-smi`.

```
nvidia-smi
```

![Saida da tela do nvidia-smi](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/14irt14klnt8fw34yv9q.png)

Se o comando não funcionar, você pode instalar o driver mais recente para sua GPU no site da [Nvidia](https://www.nvidia.com.br/Download/index.aspx?lang=br).

## Passo 2: Instalar e Configurar o Docker Desktop e WSL2

1. Se ainda não tiver o Docker Desktop instalado, baixe-o e instale-o.
2. Abra o Docker Desktop e clique em "Settings" (Configurações ⚙️).
3. Em "General", habilite a opção 'Use the WSL 2 base engine'.

![Configuração do docker desktop](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/vmplm4731m7u1praqb30.png)

## Passo 3: Instalar e Configurar o WSL2

1. No PowerShell, execute o comando para instalar a versão mais recente do WSL:

```powershell
wsl --install
```

2. Instale a distribuição do Ubuntu no WSL:

```powershell
wsl --install -d Ubuntu
```

3. Verifique se a versão do WSL está correta:

```powershell
wsl -l -v
```

![Saida do Comando wls -l -v](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/hqss6nk0ymj1v6hi4qlq.png)

Se a coluna "version" estiver como 1, atualize para a versão 2 usando o comando:

```powershell
wsl.exe --set-version NAME 2
```

Configure a versão padrão do WSL para 2 com o comando:

```powershell
wsl.exe --set-default-version 2
```

## Passo 4: Integração do Docker com WSL

No Docker Desktop, vá para "Settings" (Configurações ⚙️) > "Resources" > "WSL Integration" e marque "Enable integration with my default WSL distro".

![Tela configuração wsl integration no docker destkop](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/g465ofvwgolrkwh2yga1.png)

## Passo 5: Instalar o Toolkit do Nvidia Container

1. Remova chaves GPG antigas
```bash
sudo apt-key del 7fa2af80
```
2. Configure o repositório:

```bash
curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | sudo gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg \
&& curl -s -L https://nvidia.github.io/libnvidia-container/stable/deb/nvidia-container-toolkit.list | \
sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list
```

3. Atualize o repositório e instale o Toolkit:

```bash
sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit
```

## Passo 6: Testar a Instalação

Execute o comando `nvidia-smi` para verificar se a instalação foi bem-sucedida.

![Saáda do comando nvidia-smi no wsl ubuntu](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/iw2u8gflvi8eirpsjz39.png)

## Passo 7: Testar o Container Nvidia

Você pode testar o contêiner disponibilizado pela Nvidia usando o comando a seguir:

```powershell
docker run --gpus all nvcr.io/nvidia/k8s/cuda-sample:nbody nbody -gpu -benchmark
```

![Saida do container cuda-sample](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/kyf2cxf5p7t4xggyl86s.png)

## Extra: Testar Ambiente com Ollama

Inicie o contêiner Ollama:

```powershell
docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
```

Execute um modelo, como o TinyLlama, para verificar se está tudo funcionando corretamente:

```powershell
docker exec -it ollama ollama run TinyLlama
```

Pronto! Se tudo estiver correto, você está pronto para começar a usar o ambiente.

![Chat ollama saida](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/m9zs1b0iyc7z03wg727m.png)

## Referências:

- [Nvidia System Management Interface](https://developer.nvidia.com/nvidia-system-management-interface)
- [Docker Desktop Installation Guide](https://docs.docker.com/desktop/install/windows-install/)
- [WSL Installation Guide](https://learn.microsoft.com/en-us/windows/wsl/install)
- [Nvidia CUDA on WSL User Guide](https://docs.nvidia.com/cuda/wsl-user-guide/index.html)
- [Nvidia Container Toolkit Installation Guide](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html)
- [Ollama Documentation](https://ollama.com/)

