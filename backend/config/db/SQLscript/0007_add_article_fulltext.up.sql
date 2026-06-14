ALTER TABLE article ADD FULLTEXT INDEX ft_title_summary (title, summary) WITH PARSER ngram;
