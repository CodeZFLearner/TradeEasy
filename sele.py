from selenium import webdriver
import time

# 初始化浏览器驱动（以 Chrome 为例）
options = webdriver.ChromeOptions()
options.add_argument('--headless')  # 可选，是否以无头模式运行
options.add_argument('log-level=3')
options.page_load_strategy = 'eager'
options.binary_location = 'C:\Program Files\Google\Chrome\Application\chrome.exe'
options.add_argument("--test-type")  # 禁用沙盒模式
options.add_argument("--disable-popup-blocking")
driver = webdriver.Chrome(options=options)

try:
    # 打开目标网页
    url = "https://xueqiu.com/"  # 替换为目标 URL
    driver.get(url)
    time.sleep(1)
    # 获取当前页面的所有 Cookies
    cookies = driver.get_cookies()
    # 打印 Cookies
    for cookie in cookies:
        if cookie['name'] in ['xqat','u']:
            print(cookie['value'])
finally:
    # 关闭浏览器
    driver.quit()
