
# 交通網生成アルゴリズム

制作 : 多田 瑛貴 (ただ てるき)

![2108C70E-69D5-4492-81A7-983F2ED087EF](https://user-images.githubusercontent.com/57752033/131674598-3f096d48-1b9f-45d2-9d17-067218a0c2a3.png)


(より詳細な出力画像:[Japan.png](https://github.com/Perukii/TransportNetworkSimulator/blob/master/Japan.png))

### 概要

都市と地形図のデータをもとに、仮想的な交通網を生成するソフトです。<br>

前作(https://github.com/Perukii/TransportMaker1) と比べ、計算時間の大幅な削減を実現しています。<br>

### 設計

以下の図を参照。

![交通網生成アルゴリズム設計書](https://user-images.githubusercontent.com/57752033/117646706-08d95080-b1c7-11eb-9220-cec37656fa2f.jpg)

### 大まかなアルゴリズム

基本的な流れとしては、以下の順序で交通網の生成処理を行っています。<br>
<br>
 - 貪欲法により各都市の市街地域の範囲を推測<br>
 - 市街地域の繋がっている複数の都市を一つの都市圏としてまとめ上げる<br>
 - 対象の都市を結ぶ最小全域木を構築、ルートを推測<br>
 - A\*探索(最適経路探索)により各ルートのより細かなパスを算出<br>
 - 一部の例外的なパスの消去(海の上を通り過ぎるパスなど)<br>
<br>
※実際はさらに細かな処理を行っています。<br>
<br>

## データ出典元・参考文献

### 使用データ

地理院地図<br>
https://maps.gsi.go.jp/#5/36.261992/137.285156/&ls=relief_free&disp=1&lcd=relief_free&vs=c1j0h0k0l0u0t0z0r0s0m0f1&reliefdata=20G00FF00G2FDG000000G5FAG0000FFG8F7GFF00FFGGFF00FF<br>

アマノ技研 地方公共団体の位置データ<br>
https://amano-tec.com/data/localgovernments.html<br>

e-Stat 統計で見る日本<br>
https://www.e-stat.go.jp/regional-statistics/ssdsview/municipality<br>

### 参考文献

地図における数学<br>
https://w3e.kanazawa-it.ac.jp/e-scimath/contents/t15/textbook_t15_all.pdf<br>

