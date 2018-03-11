# piyopoked

常駐式のぴよポケワンドロボット  

**piyopoke + d(aemon)**  

## 環境構築メモ

ドコモのAPI取った
https://dev.smt.docomo.ne.jp/?p=docs.api.index

prebuild
```
go install github.com/Mushus/piyopoked/vendor/layeh.com/gopus
```

ffmpegインストール(WSLだと一番下以外不要だった)
```
sudo apt-get -y install software-properties-common
sudo add-apt-repository ppa:mc3man/trusty-media
sudo apt-get update
sudo apt-get -y install ffmpeg
```

youtube-dlインストール(WSLだとシンボリックリンクが必要だった)
```
sudo wget https://yt-dl.org/downloads/latest/youtube-dl -O /usr/local/bin/youtube-dl
sudo chmod a+rx /usr/local/bin/youtube-dl
sudo ln -s /usr/bin/python3 /usr/bin/python
```

コーデックいる？
```
sudo apt-get install libav-tools
```

openjtalkインストール
```
sudo apt-get -y install open-jtalk open-jtalk-mecab-naist-jdic hts-voice-nitech-jp-atr503-m001
```

mei install
```
wget http://downloads.sourceforge.net/project/mmdagent/MMDAgent_Example/MMDAgent_Example-1.7/MMDAgent_Example-1.7.zip
unzip MMDAgent_Example-1.7.zip
sudo cp -R ./MMDAgent_Example-1.7/Voice/mei /usr/share/hts-voice/
```

招待
```
https://discordapp.com/oauth2/authorize?client_id=YOUR_CLIENT_ID&scope=bot&permissions=0
```
