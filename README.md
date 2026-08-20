# gear-contact

`gear-contact` 是一个渐开线直齿圆柱齿轮（外啮合）啮合核算命令行工具。输入一个 JSON 文件给出模数 m、两齿轮齿数 z1/z2、压力角 α（度）、齿顶高系数 ha*（可覆盖，默认 1）、顶隙系数 c*（默认 0.25）与可选变位系数 x1/x2；输出两轮的分度圆 d=m·z、基圆 db=d·cosα、齿顶圆 da、中心距、法向齿距、作用线长度、重合度 εα、顶隙与根切状态。边界：m≤0、齿数非正整数、齿数过小且未设足够变位导致根切、压力角超出钉定区间（14.5°~25°，且必须 <90°）均为错误，打印到 stderr 并以非零码退出。重合度不做查表，一律按「作用线长度 ÷ 法向齿距」得出。

## 用法

```bash
gear-contact mesh example/spur-20.json
gear-contact help
```

算例 `example/spur-20.json`：α=20°、z1=20、z2=40、m=1，重合度应 >1（约 1.635）。

```
gear-contact mesh
module m = 1.000 mm
alpha = 20.000 deg, alpha' = 20.000 deg (standard)
gear 1: z=20 x=+0 d=20.000 db=18.794 da=22.000 df=17.500
gear 2: z=40 x=+0 d=40.000 db=37.588 da=42.000 df=37.500
centre distance a = 30.000 mm (standard 30.000 mm)
base pitch pb = 2.952 mm
line of action length = 4.827 mm
contact ratio epsilon_alpha = 1.635
tip clearance c = 0.250 mm (positive)
undercut: none
```

## 输入 JSON

```json
{
  "module": 1.0,
  "pressure_angle_deg": 20.0,
  "addendum_coefficient": 1.0,
  "clearance_coefficient": 0.25,
  "gear1": { "teeth": 20, "shift": 0.0 },
  "gear2": { "teeth": 40, "shift": 0.0 }
}
```

字段说明：

| 字段 | 含义 | 默认/约束 |
| --- | --- | --- |
| `module` | 模数 m（mm） | 必须 >0 |
| `pressure_angle_deg` | 标准压力角 α（度） | 钉定 14.5°~25°，<90° |
| `addendum_coefficient` | 齿顶高系数 ha* | 默认 1 |
| `clearance_coefficient` | 顶隙系数 c* | 默认 0.25 |
| `gear1` / `gear2` | `teeth` 齿数（正整数）、`shift` 变位系数 x | shift 默认 0 |

未知字段会被拒绝（防拼写错误静默改变几何）。两轮共用同一模数与压力角，不存在主从动各自一套不一致参数。

## 几何公式（钉定）

- 分度圆 `d = m·z`，基圆 `db = d·cosα`
- 齿顶圆 `da = d + 2·m·(ha* + x)`（标准 ha*=1）
- 齿根圆 `df = d − 2·m·(ha* + c* − x)`
- 标准中心距 `a0 = m·(z1+z2)/2`，标准安装时等于两分度圆半径之和
- 变位：`x1+x2=0` 时工作压力角 `α' = α`、工作中心距 `a = a0`；否则解无侧隙方程
  `inv(α') = inv(α) + 2·(x1+x2)·tanα/(z1+z2)`（`inv(x)=tan x−x`，牛顿迭代求反函数），
  再按基圆不变量 `a·cosα' = a0·cosα` 得 `a`
- 法向齿距 `pb = π·m·cosα`（**必须乘 cosα**）
- 作用线：两基圆内公切线，与两齿顶圆的交点之间为啮合线段，长度
  `gα = √(ra1²−rb1²) + √(ra2²−rb2²) − a·sinα'`
- 重合度 `εα = gα / pb`
- 顶隙 `c = a − ra1 − rf2`（对称地 `= a − ra2 − rf1`），标准齿高下为正且 `c = c*·m = 0.25·m`
- 根切：最小齿数 `zmin = 2·ha*/sin²α`，最小变位 `xmin = ha*·(1 − z/zmin)`
  （α=20°、ha*=1 时即经典的 `xmin ≈ (17−z)/17`）。`x < xmin` 判定根切，**钉为 error**：
  校验时直接拒绝该齿轮副（`m=0`、`z≤0`、`α≥90°`、未设变位的过小齿数均同）

## 结构

- `main.go` 只做子命令接线与 JSON 读取，求解全部在 `internal/`
- `internal/angle`：度弧度换算、压力角区间校验、inv 函数与反函数求解
- `internal/gears`：齿轮轮廓（四个特征圆）、JSON 规格、校验与根切判定、齿轮副
- `internal/mesh`：中心距、工作压力角、作用线端点与长度、重合度、顶隙、报告

## 构建与测试

```bash
go build ./...
go test ./...
go test -run TestExampleContactRatioAboveOne ./internal/mesh/
```

## 交叉规则

- 互换主从动齿数：中心距与重合度不变
- 很大齿数配对仍 `εα > 1`
- 标准安装中心距 = 两分度圆半径之和；标准齿顶高时顶隙为正

## 许可证

MIT，见 `LICENSE`。
