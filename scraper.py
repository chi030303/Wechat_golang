from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time
import sys
import json
from selenium.common.exceptions import NoSuchElementException

sys.stdout.reconfigure(encoding='utf-8')  

# 爬取数据的列表
data_list = []
# 爬取数据的id
idx = 0

def get_info(driver):
    try:
    # 定位内容1和内容2的元素
        content1_element = driver.find_element(By.XPATH, '//*[@id="app"]/div/div[1]/div[2]/div/div[3]/div/div[1]/div[1]/div/div[1]/div[1]')
        content2_element = driver.find_element(By.XPATH, '//*[@id="app"]/div/div[1]/div[2]/div/div[3]/div/div[1]/div[1]/div/div[2]/div/div[1]/span')
    except NoSuchElementException as e:
        print("找不到元素:", e)
        return None
    # 提取内容1和内容2的文本内容
    content1_text = content1_element.text
    content2_text = content2_element.text
    
    recruitment_target = ""
    recruitment_position = ""
    recruitment_major = ""
    recruitment_date = ""
    company_name = ""
    recruitment_batch = ""
    job_city = ""
    recruitment_session = ""

    try:
        company_element = driver.find_element(By.XPATH, '//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[1]/p')
        company_name = company_element.text
    except NoSuchElementException as e:
        print("找不到公司名称元素:", e)

    try:
        # 定位内容2下面的其他元素
        other_elements = driver.find_elements(By.XPATH, '//*[@id="app"]/div/div[1]/div[2]/div/div[3]/div/div[1]/div[1]/div/div[2]/div/div[position()>1]/span')
        elements = [element.text for element in other_elements]
    except NoSuchElementException as e:
        print("找不到其他元素:", e)

    # 将爬取到的文本处理成列表
    lines = content1_text.split("\n")

    # 对列表的内容进行筛选并存储
    for i in range(len(lines)):
        if i < len(lines) and (lines[i]) == "招聘对象：":
            i += 1
            while i < len(lines) and "：" not in lines[i]:
                recruitment_target += lines[i]
                i += 1
        if i < len(lines) and (lines[i]) == "招聘岗位：":
            i += 1
            while i < len(lines) and "：" not in lines[i]:
                recruitment_position += lines[i]
                i += 1
        if i < len(lines) and ((lines[i]) == "招聘专业：" or "招聘要求：" in (lines[i]) or "招聘条件：" in (lines[i])):
            i += 1
            while i < len(lines) and "：" not in lines[i]:
                recruitment_major += lines[i]
                i += 1

    recruitment_batch += content2_text
    
    job_city = elements[0]
    recruitment_date = elements[1]
    recruitment_session = elements[-1]
    # print(f"公司名称：{company_name}，招聘对象：{recruitment_target},招聘岗位：{recruitment_position},招聘专业/条件/要求：{recruitment_major},工作城市：{job_city}，投递时间：{recruitment_date}，招聘届数：{recruitment_session}")
    
    # 声明idx是全局变量
    global idx
    idx += 1
    data = {
        "id": idx,
        "公司名称": company_name,
        "招聘对象": recruitment_target,
        "招聘岗位": recruitment_position,
        "招聘专业/条件/要求": recruitment_major,
        "工作城市": job_city,
        "投递时间": recruitment_date,
        "招聘届数": recruitment_session
    }
    
    # 将每个招聘信息字典添加到列表中
    data_list.append(data)
    
# 初始化 WebDriver
driver = webdriver.Edge()

# 打开网页
driver.get("https://offershow.cn/jobs/specialarea")

# 爬取网页信息
def crawl_pages():
    # 定义一个循环来翻页，共106页
    for page in range(104):
        time.sleep(1)
        # 在每一页循环点击每个窗口
        for window in range(1, 11):
            # 构造窗口的 XPath
            window_xpath = f'//*[@id="app"]/div/div[2]/div[2]/div[2]/div[1]/div[2]/div[{window}]'
            # 点击窗口
            try:
                WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.XPATH, window_xpath)))
                element = WebDriverWait(driver, 10).until(EC.element_to_be_clickable((By.XPATH, window_xpath)))
                element.click()
            except Exception as e:
                print("点击窗口出错:", e)
                continue  # 如果出现异常，跳过当前循环，继续下一次循环

            # 等待新窗口打开
            WebDriverWait(driver, 10).until(EC.number_of_windows_to_be(2))

            # 切换到新窗口
            driver.switch_to.window(driver.window_handles[-1])

            # 等待新窗口加载完成
            time.sleep(1)

            # 在这里爬取你需要的信息
            get_info(driver)
            # print(c1, c2, c3)

            # 关闭当前窗口
            driver.close()

            # 切换回主窗口
            driver.switch_to.window(driver.window_handles[0])

            # 将爬取到的信息追加到json文件末尾
            with open('companies_data.json', 'w', encoding='utf-8') as f:
                json.dump(data_list, f, ensure_ascii=False, indent=4)

        # 模拟点击下一页按钮
        try:
            next_button = driver.find_element(By.XPATH,'//*[@id="page_change"]/div/button[2]')
            next_button.click()
        except Exception as e:
            print("点击下一页按钮出错:", e)
            return None

if __name__ == "__main__":

    # 爬取页面信息
    crawl_pages()

    # 关闭浏览器
    driver.quit()
