#!/usr/bin/env python
# -*- coding: utf-8 -*-

from graphviz import Digraph

# formatはpngを指定(他にはPDF, PNG, SVGなどが指定可)
G = Digraph(format='png')

C = Digraph(name='creator', format='png')
C.attr('node', shape='circle', style='filled')
C.attr('graph', rank='source')

P = Digraph(name='peer', format='png')
P.attr('node', shape='trapezium', style='filled', fillcolor='#FFEECC')
P.attr('graph', rank='sink')

N = 10      # Creator数
M = 4       # Peer数

def cr(a):
    return "C"+str(a)

def pr(a):
    return "P"+str(a)

# ノードの追加
for i in range(N):
    C.node(cr(i), cr(i))

for i in range(M):
    if i == 0:
        P.node(pr(i), 'root')
    else:
        P.node(pr(i), pr(i))

G.subgraph(C)
G.subgraph(P)

for i in range(N):
    G.edge(cr(i), pr(i%M), color='#FF8822', weight="0.1")

# binary_tree.pngで保存
G.render('_demo_graph01')

# 辺の追加
for i in range(N):
    G.edge(cr(i), cr((i+1)%N))
    if i < 8:
        G.edge(cr(i), cr((i+2)%N))
    if i < 6:
        G.edge(cr(i), cr((i+3)%N))
    if i < 4:
        G.edge(cr(i), cr((i+4)%N))

    if i > 5:
        G.node(cr(i), label=cr(i), style='filled', fillcolor='#FFCC33')

for i in range(M):
    G.node(pr(i), pr(i))

G.render('_demo_graph02')

# 辺の追加
for i in range(N):
    G.edge(cr(i), cr(2))




