from selenium import webdriver
from bs4 import BeautifulSoup
import urllib2
import sys
import re
import os

start = 59
def main():
	try:
		os.mkdir("./downloadedMusic")
	except:
		print("Who cares")
	os.chdir("./downloadedMusic")
	if sys.argv[1] == "a":
		artistDownload(sys.argv[2])
	if sys.argv[1] == "g":
		genreDownload(sys.argv[2])

def genreDownload(genre):
	num = 0
	sent = False
	for i in range(0, 79) :
			if i == 0:
				response = urllib2.urlopen('http://www.midiworld.com/search/?q='+genre)
				print('http://www.midiworld.com/search/?q='+genre)
			else:
				response = urllib2.urlopen('http://www.midiworld.com/search/' + str(i) +'/?q='+genre)
				print('http://www.midiworld.com/search/' + str(i) +'/?q='+genre)
				html = response.read()
				for match in re.findall(r'([A-Z].*) - <a href="(.*)" target.*>download</a>', html) :
					print("Downloading midi file: %s"%match[0])
					os.system('wget %s --output-document=\"%s-%s.mid\"'%(match[1], num, match[0].replace(" ","_")))
					print('-'*20)
					num += 1
					if sent:
						break
	"""
	browser = webdriver.Chrome()
	for i in range(start,100000):
		#try:
			url ='https://freemidi.org/download-'+str(i)

			u = urllib2.urlopen(url)
			html = u.read().decode('utf-8')
			soup = BeautifulSoup(html)

			for link in soup.select('ol'):
				print(link.GetText())
				passed =False
				if genre in link.getText() and not passed:
					browser.get(url)
					button =browser.find_element_by_id('downloadmidi')
					button.click()
					print('pass')
					passed= True
					print(i)
				print(i)
		#except Exception as e:
			#print(e)
			#continue
			"""
def artistDownload(genre):
	num = 0
	sent = False
	for i in range(0, 79) :
			if i == 0:
				response = urllib2.urlopen('http://www.midiworld.com/search/?q='+genre)
				print('http://www.midiworld.com/search/?q='+genre)
			else:
				response = urllib2.urlopen('http://www.midiworld.com/search/' + str(i) +'/?q='+genre)
				print('http://www.midiworld.com/search/' + str(i) +'/?q='+genre)
				html = response.read()
				for match in re.findall(r'([A-Z].*) - <a href="(.*)" target.*>download</a>', html) :
					print("Downloading midi file: %s"%match[0])
					os.system('wget %s --output-document=\"%s-%s.mid\"'%(match[1], num, match[0].replace(" ","_")))
					print('-'*20)
					num += 1
					if sent:
						break
main()
