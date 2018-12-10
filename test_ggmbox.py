import os
import re
import shutil
import tempfile
import unittest

from scrapy.crawler import CrawlerProcess
from scrapy.utils.log import configure_logging
from twisted.internet import reactor, task

from ggmbox import GoogleGroupMBoxSpider


class GGMBoxTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        configure_logging()

    def setUp(self):
        self.process = CrawlerProcess()
        self.tempdir = tempfile.mkdtemp(prefix="ggmbox-")
        self.process.crawl(GoogleGroupMBoxSpider, name="golang-nuts", output=self.tempdir)

    def tearDown(self):
        shutil.rmtree(self.tempdir)

    def test_integration(self):
        task.deferLater(reactor, 8, reactor.stop)
        reactor.run()
        dirs = os.listdir(self.tempdir)
        self.assertGreater(len(dirs), 0)
        files_fetched = False
        for dirname in dirs:
            files = os.listdir(os.path.join(self.tempdir, dirname))
            if len(files) > 0:
                files_fetched = True
                files.sort()
                self.assertGreater(len(files), 0)
                re.match(r"\d{3}_", files[0][:4])
                self.assertEqual(files[0][-6:], ".email")
                self.assertGreater(
                    os.path.getsize(os.path.join(self.tempdir, dirname, files[0])), 0)
        self.assertTrue(files_fetched)


if __name__ == "__main__":
    unittest.main()
