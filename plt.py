import json
import matplotlib.pyplot as plt
from matplotlib.font_manager import FontProperties

# 设置中文字体
font_path = '/System/Library/Fonts/PingFang.ttc'  # macOS系统字体
font_prop = FontProperties(fname=font_path)

# 读取 JSON 数据
with open('/tmp/rank.json', 'r') as f:
    data = json.load(f)

# 提取总成绩、排名和面试成绩数据
total_scores = [result['total_score'] for result in data['results']]
your_ranks = [result['your_rank'] for result in data['results']]
your_interview_scores = [result['rankings'][-1]['interview_score'] for result in data['results']]
top5_probability = data['top5_probability']

# 提取每次模拟中的最高和最低总成绩
max_scores = [max(result['rankings'], key=lambda x: x['total_score'])['total_score'] for result in data['results']]
min_scores = [min(result['rankings'], key=lambda x: x['total_score'])['total_score'] for result in data['results']]

# 绘制图表
plt.figure(figsize=(12, 8))

# 绘制总成绩曲线，包括最高和最低范围
plt.subplot(3, 1, 1)
plt.plot(total_scores, label='总成绩')
plt.plot(max_scores, label='最高总成绩', linestyle='--')
plt.plot(min_scores, label='最低总成绩', linestyle='--')
plt.fill_between(range(len(total_scores)), min_scores, max_scores, color='gray', alpha=0.2)
plt.xlabel('模拟次数', fontproperties=font_prop)
plt.ylabel('总成绩', fontproperties=font_prop)
plt.title('模拟中的总成绩变化', fontproperties=font_prop)
plt.legend(prop=font_prop)

# 绘制面试成绩曲线
plt.subplot(3, 1, 2)
plt.plot(your_interview_scores, label='你的面试成绩', color='green')
plt.xlabel('模拟次数', fontproperties=font_prop)
plt.ylabel('面试成绩', fontproperties=font_prop)
plt.title('模拟中的面试成绩变化', fontproperties=font_prop)
plt.legend(prop=font_prop)

# 绘制排名曲线
plt.subplot(3, 1, 3)
plt.plot(your_ranks, label='你的排名', color='orange')
plt.axhline(y=5, color='r', linestyle='--', label='前五名阈值')
plt.xlabel('模拟次数', fontproperties=font_prop)
plt.ylabel('你的排名', fontproperties=font_prop)
plt.title(f'模拟中的排名变化 (进入前五名的概率: {top5_probability:.2%})', fontproperties=font_prop)
plt.legend(prop=font_prop)

plt.tight_layout()
plt.show()
