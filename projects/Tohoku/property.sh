
HEIGHT_DIFFELENCE="3"
IMAGE_WIDTH="2075"
IMAGE_HEIGHT="1845"
LATITUDE_START="37.7"
LATITUDE_END="41.6"
LONGITUDE_START="137.8"
LONGITUDE_END="143.5"

# スコア計算における高さの影響度
HEIGHT_SCORE="0.1"
# スコア計算における高低差の影響度
HEIGHT_DIFFELENCE_SCORE="1.4"
# スコア計算における2都市間の距離の影響度
DISTANCE_SCORE="1800.0"
# スコア計算における都市圏範囲の影響度
URBAN_AREA_SCORE="0.0"
# スコア計算における海の影響度
SEA_AREA_SCORE="2000.0"
# スコア計算における人口の影響度
POPULATION_SCORE="0.0000006"
# クラスカル路を適用する都市の最小の人口
KRUSKAL_PATH_MIN_POPULATRION="15000"
# 非クラスカル路を採用するのに満たす必要のある、一つの都市に対する対象のパスと他都市間のパスとの角度の差の最小値の下限
KRUSKAL_PATH_MAX_ANGLE_DIFFELENCE="80.0"
# パスのスコア計算にて、一つの都市が比較を行う都市の数
KRUSKAL_PATH_MAX_CROSS="8"
# パスの生成に利用するA*路の粗さ
PATH_RELEASE_INTERVAL="0.005"
# パスのコスト計算に利用するA*路の粗さ
PATH_DRAFT_INTERVAL="0.02"
# 都市圏の範囲の定義に利用するA*路の粗さ
URBAN_AREA_INTERVAL="0.0022"
# 広域都市圏の範囲の密集度
URBAN_WIDE_AREA_DENSITY="50"
# 都市圏の範囲の密集度
URBAN_CENTRAL_AREA_DENSITY="200"
