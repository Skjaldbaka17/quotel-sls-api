SELECT 
*, 
ts_rank(quote_tsv, plainq) as plainrank, 
ts_rank(quote_tsv, phraseq) as phraserank, 
ts_rank(quote_tsv, generalq) as generalrank 
FROM 
searchview, 
plainto_tsquery('float like butterfly') as plainq, 
to_tsquery('float <-> like <-> butterfly') as phraseq,
to_tsquery('float | like | butterfly') as generalq  
WHERE 
( 
    tsv @@ plainq OR 
    tsv @@ phraseq OR 
    'float like butterfly' % ANY(STRING_TO_ARRAY(name,' ')) OR 
    tsv @@ generalq
) 
ORDER BY 
phraserank DESC,
similarity(name, 'float like butterfly') DESC, 
plainrank DESC, 
generalrank DESC, 
author_id DESC 
LIMIT 25;

SELECT *, 
ts_rank(tsv, plainq) as plainrank 
FROM quotes, 
plainto_tsquery('float like butterfly') as plainq  
WHERE ( tsv @@ plainq ) 
ORDER BY plainrank desc, 
author_id desc LIMIT 25;


select *, similarity(name,'friedrik nietse') as sml from authors where name % 'friedrik nietse' order by sml desc limit 10;
select *, similarity(quote,'float like a butterfly') as sml from quotes where quote % 'friedrik nietse' order by sml desc limit 10;
CREATE INDEX quotes_trgm_idx ON quotes USING GIN (quote gin_trgm_ops);