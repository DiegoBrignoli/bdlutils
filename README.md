bdlpdf
===
Bdlpdf is a simple API written in GO and PHP that generates a PDF document from the site https://www.bdl.servizirl.it using the document ID.

EX:

https://www.bdl.servizirl.it/bdl/bookreader/index.html?path=fe&cdOggetto= **17014**

<h3>Methods</h3>

You can only use the method via **GET**

<h3>Request syntax</h3>

http://waifuai.dd-dns.de:801/bdlpdf/index.php
 ? [id=<document id>]
 & [pag=<number of page>]

<h3>Parameters</h3>

* **Id**: The id of the document obtainable on the site https://www.bdl.servizirl.it
  * Id > 0
  * Obligatory parameter
* **Pag**: The number of pages to be obtained and inserted into the pdf document
  * **Default**: 'all'
  * pag > 0
  * pag = 'all': Gets and inserts all pages in the pdf
  * Not obligatory parameter

<h3>Example</h3>

1. http://waifuai.dd-dns.de:801/bdlpdf/index.php?id=17014
2. http://waifuai.dd-dns.de:801/bdlpdf/index.php?id=17014&pag=all
3. http://waifuai.dd-dns.de:801/bdlpdf/index.php?id=17014&pag=3
