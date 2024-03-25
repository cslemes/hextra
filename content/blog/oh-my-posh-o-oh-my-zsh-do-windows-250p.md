---
title: '"Oh My Posh: O Oh My Zsh do Windows"'
author: '"https://media.dev.to/cdn-cgi/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Fmmxr5clqgi86kjbqetgo.png"'
date: '"2024-02-15T23:09:29.755Z"'
coverImage: ' "https://media.dev.to/cdn-cgi/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Fmmxr5clqgi86kjbqetgo.png"'
ogImage: '"https://media.dev.to/cdn-cgi/image/width=1000,height=420,fit=cover,gravity=auto,format=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2Fmmxr5clqgi86kjbqetgo.png"'
excerpt: '"O Oh My Zsh é uma ferramenta bem conhecida para customizar o terminal Zsh, usado em distribuições..."'
---

O Oh My Zsh é uma ferramenta bem conhecida para customizar o terminal Zsh, usado em distribuições Linux e no macOS. Além de oferecer uma variedade de temas, o Oh My Zsh também possui uma ampla gama de plugins, como completions, que ajudam a aumentar a produtividade.

Para aqueles que utilizam o Windows e desejam explorar algumas possibilidades oferecidas pelo Oh My Zsh, decidi escrever este artigo sobre o Oh My Posh. O Oh My Posh pode ser usado em vários shells, não apenas no PowerShell, e também é multiplataforma. Isso significa que ele pode ser executado tanto no Windows quanto em qualquer shell do Linux e no macOS, além de ser gratuito e de código aberto. No entanto, ele se limita à funcionalidade de temas e não oferece gerenciamento de plugins como o Oh My Zsh. O PowerShell, por sua vez, possui várias funções de completions que podem ser configuradas, mas estão fora do escopo deste artigo. Demonstrarei como instalá-lo no Windows.

O primeiro passo é garantir que você tenha um terminal adequado no Windows. Para isso, vá até a Microsoft Store e baixe o Windows Terminal. Você também pode usar um gerenciador de pacotes via linha de comando.

- Via Winget (nativo do Windows 11):

```powershell
winget install --id Microsoft.WindowsTerminal -e
```

- Via [Chocolatey](https://chocolatey.org/) (não oficial):

```powershell
choco install microsoft-windows-terminal
```

- Via [Scoop](https://scoop.sh/) (não oficial):

```powershell
scoop bucket add extras
scoop install windows-terminal
```

Em seguida, instale o PowerShell Core. No Windows padrão, o PowerShell antigo vem pré-instalado, mas você pode optar por instalar o PowerShell Core usando o pacote .msi disponibilizado pela Microsoft ou usando o Winget. Como não encontrei referências na documentação da Microsoft sobre a instalação usando outros gerenciadores, não posso recomendar neste momento.

- Via MSI:

[PowerShell-7.4.1-win-x64.msi](https://github.com/PowerShell/PowerShell/releases/download/v7.4.1/PowerShell-7.4.1-win-x64.msi)

- Via Winget:

```

winget install --id Microsoft.PowerShell --source winget

```

Após instalar o PowerShell Core, abra o Windows Terminal e defina-o como o terminal padrão. Para fazer isso, vá em Configurações, clique na seta para baixo ao lado do botão '+' na aba de título da janela e, em seguida, em Perfil Padrão, escolha PowerShell (não Windows PowerShell).

Agora, finalmente, vamos instalar o Oh My Posh, seguindo as instruções da documentação oficial. Existem opções para instalar usando gerenciadores de pacotes, mas neste caso, vamos utilizar o script PowerShell para instalar.

No Windows Terminal com PowerShell Core, execute a seguinte linha de comando e pressione Enter:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://ohmyposh.dev/install.ps1'))

```

Agora, você deve configurar o Oh My Posh para iniciar junto com seu terminal. Para isso, edite o arquivo de perfil do PowerShell. A localização do arquivo fica armazenada na variável de ambiente **`$PROFILE`**. Você pode usar o Notepad para isso:

```
notepad $PROFILE
```

Vá para a última linha do arquivo e adicione:

```
oh-my-posh init --shell pwsh | Invoke-Expression
```

Salve e feche o Notepad, e reinicie seu terminal. Deve abrir conforme a imagem abaixo, porém sem os ícones.

![Tela do powershell com retangulos ao no lugar do caracteres](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/v5d07mrguh6jnja5ndok.png)

Para carregar os ícones, é necessário ter uma fonte que suporte esses caracteres. Para isso, vamos utilizar o [Nerd Fonts](https://www.nerdfonts.com/). Basta escolher uma fonte, baixá-la e instalá-la com um duplo clique no Windows. Você também pode fazer isso por linha de comando usando o cmdlet do Oh My Posh para instalar fontes. Por exemplo, para instalar a fonte MesloLG, você pode executar o seguinte comando:

```powershell
oh-my-posh font install MesloLG
```

Agora, é necessário ir nas configurações do Windows Terminal. Em "Perfis" > "Padrões" > "Aparência", selecione o tipo de fonte e troque para a fonte escolhida. Salve e feche as configurações. O Windows Terminal aplicará as alterações na janela atual.

Referências

https://learn.microsoft.com/pt-br/windows/terminal/
https://learn.microsoft.com/pt-br/powershell/scripting/install/installing-powershell-on-windows
https://ohmyposh.dev/
