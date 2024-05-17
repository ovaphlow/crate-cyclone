# 使用官方的Python基础镜像
FROM python:3.12-alpine3.19

# 设置工作目录
WORKDIR /app

# 将当前目录的内容复制到工作目录中
COPY . /app

# 安装项目依赖
RUN pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --no-cache-dir -r requirements.txt

# 设置环境变量
ENV NAME World

# 运行应用
CMD ["uvicorn", "main:app", "--reload", "--host", "0.0.0.0", "--port", "8421"]
