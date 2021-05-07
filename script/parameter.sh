# スコア計算における高さの影響度
HEIGHT_SCORE="0.1"
# スコア計算における高低差の影響度
HEIGHT_DIFFELENCE_SCORE="1.0"
# スコア計算における2都市間の距離の影響度
DISTANCE_SCORE="2500.0"
# スコア計算における都市圏範囲の影響度
URBAN_AREA_SCORE="0.0"
# スコア計算における海の影響度
SEA_AREA_SCORE="10000.0"
# スコア計算における人口の影響度
POPULATION_SCORE="0.0"
# クラスカル路を適用する都市の最小の人口
KRUSKAL_PATH_MIN_POPULATRION="10000"
# 非クラスカル路を採用するのに満たす必要のある、一つの都市に対する対象のパスと他都市間のパスとの角度の差の最小値の下限
KRUSKAL_PATH_MAX_ANGLE_DIFFELENCE="70.0"
# パスのスコア計算にて、一つの都市が比較を行う都市の数
KRUSKAL_PATH_MAX_CROSS="10"
# パスの生成に利用するA*路の粗さ
PATH_RELEASE_INTERVAL="0.005"
# パスのコスト計算に利用するA*路の粗さ
PATH_DRAFT_INTERVAL="0.008"
# 都市圏の範囲の定義における高低差の影響度
URBAN_AREA_HEIGHT_DIFFELENCE_SCORE="4.0"
# 広域都市圏の範囲の密集度
URBAN_WIDE_AREA_DENSITY="80"
# 都市圏の範囲の密集度
URBAN_CENTRAL_AREA_DENSITY="300"
# 直線距離あたりのパスの距離の最大
MAX_PATH_DISTANCE_PER_CITY_DISTANCE="1.6"
# 海を渡るパスの距離の最大
MAX_BRIDGE_DISTANCE="0.1"
# 都市と都市の間の直線距離の最大(2個目以降のパスに適用)
MAX_CITY_DISTANCE="0.8"