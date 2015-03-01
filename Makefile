BASE=github.com/hayeah/jsonql

sample.db:
	sqlite3 sample.db < sampledb.sql

install:
	go install $(BASE)
	go install $(BASE)/fixture
	go install $(BASE)/handler
	go install $(BASE)/jsonql