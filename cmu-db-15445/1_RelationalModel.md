# Intro and Relational Model

## Concepts

- Data Models
    A data model is collection of concepts for describing the data in a database
    + Relational -> the most DBMSs
    + Key/Value -> nosql
    + Graph -> nosql
    + Document -> nosql
    + Column-Family -> nosql
    + Array/Matrix -> Machine Learning
    + Hierarchical, rare
    + Network, rare
- Schema
    + A schema is a description of a particular collection of data, using a given data model


## Relational Model
- Primary keys  
    A relation's primary key uniquely identifies a single tupe.
- Foreign keys  
    Specify that an attribute from one relation has map to a tuple in another relation.
- Relational algebra
    + σ Select
    + π Projection
    + u Union
    + n Intersection
    + - Difference
    + x Product
        * Generate a relation that contains all possible combinations of tuples from the input relations.
    + ⋈ Join
        * combinations of two tuples with a common values for one or more attributes.
- Queries  
    The relational model is independent of any query language implementation. 
