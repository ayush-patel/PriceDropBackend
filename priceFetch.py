import sys
import requests
from bs4 import BeautifulSoup
import html5lib

def priceFetch(url):
    try:
        webpage = requests.get(url)
        price = BeautifulSoup(webpage.text, 'html5lib').body.find("div", class_= 'current-price').text
        return price[1:]
    except:
        return "Error: Couldn't fetch the price"

if __name__ == "__main__":
    print priceFetch("")
