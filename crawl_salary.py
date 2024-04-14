from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time
import sys
import json
from selenium.common.exceptions import NoSuchElementException

sys.stdout.reconfigure(encoding='utf-8')  

salary_list = []
idx = 0 

def crawl_salaries():
    url = "https://offershow.cn/?search_tab=2"
    driver = webdriver.Edge()
    driver.get(url)
    for page in range(21):
        time.sleep(4)
        for i in range(1, 22):
            try:
                if is_div1_exists(driver, i):
                    job_position = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[1]/p').text
                    salary = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[1]/span').text
                    work_location = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[2]/span[1]').text
                    education_background = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[2]/span[2]').text
                    major = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[2]/span[3]').text
                    company_name = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[2]/div/p').text
                    company_profile = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{i}]/div[2]/span').text
                else:
                    job_position = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[1]/p').text
                    salary = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[1]/span').text
                    work_location = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[2]/span[1]').text
                    education_background = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[2]/span[2]').text
                    major = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[1]/div[2]/span[3]').text
                    company_name = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[2]/div/p').text
                    company_profile = driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div/div[2]/div/div[4]/div[1]/div[{i}]/div[2]/span').text
            except NoSuchElementException as e:
                print("找不到元素:", e)
                continue

                    # 声明idx是全局变量
            global idx
            idx += 1
            data = {
                "id": idx,
                "公司名称": company_name,
                "薪资": salary,
                "工作城市": work_location,
                "学历要求": education_background,
                "专业": major,
                "岗位": job_position,
                "公司简介": company_profile
            }
            
            # 将每个招聘信息字典添加到列表中
            salary_list.append(data)
            print(f"岗位：{job_position}，薪资：{salary}，城市：{work_location}，学历：{education_background}，专业：{major}，公司名称：{company_name}，公司简介：{company_profile}")
        
        # 将爬取到的信息追加到json文件末尾
        with open('salaries_data.json', 'w', encoding='utf-8') as f:
            json.dump(salary_list, f, ensure_ascii=False, indent=4)

        # 模拟点击下一页按钮
        try:
            next_button = driver.find_element(By.XPATH,'//*[@id="page_change"]/div/button[2]')
            next_button.click()
        except Exception as e:
            print("点击下一页按钮出错:", e)
            continue

# 由于xpath有所不同，所有动态地进行调整
def is_div1_exists(driver, index):
    try:
        driver.find_element(By.XPATH, f'//*[@id="app"]/div/div[1]/div[2]/div[1]/div[2]/div/div[4]/div[1]/div[{index}]/div[1]')
        return True
    except NoSuchElementException:
        return False

if __name__ == "__main__":
    crawl_salaries()