import pickle
import sys

PATH = "/home/rohan/Documents/Lusio/API/"
with open(PATH+'scores.pkl', 'rb') as f:
    scores = pickle.load(file=f)

i = str(sys.argv[1])
k = int(sys.argv[2])


r = [a for a in sorted(scores[240], key=lambda t: t[1], reverse=True)]
    
    
i = int(i[1:])-1
ids = [str('T'+'0'*(3-len(str(a[0]+1)))+str(a[0]+1)) for a in sorted(scores[i], key=lambda t: t[1], reverse=True)[:k]]

print(ids)

