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
author_id desc LIMIT 25