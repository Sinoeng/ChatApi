import os
import mysql.connector as connector

class Python_db:
    db
    def __init__(self):
        print("Connecting to database")
        self.__db = connector.connect(
            host=os.environ['DB_HOST'],
            user="bob",
            password="bob"
        )
        print("connected to database")

    def test_exec(self):
        ans = self.__db.cursor().execute("SHOW DATABASES;")
        print("Answer is", ans)



