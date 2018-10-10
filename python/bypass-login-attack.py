
# remove put
# FF 15 44 10 40 00 
# replace with NOPs at btye 1035
# 90 90 90 90 90 90 
replacement_hex_1 = '909090909090'

# remove scanf
# FF 15 4C 10 40 00
# replace with NOPs at btye 1053
# 90 90 90 90 90 90
replacement_hex_2 = '909090909090'

# remove cmp
# 3B 05 6C 30 40 00
# replace with NOPs at byte 1076    
# 90 90 90 90 90 90
replacement_hex_3 = '909090909090'

# remove jnz
# 75 1E
# replace with NOPs at byte 1082
# 90 90 
replacement_hex_4 = '9090'

import re
import os
import sys

cmd = 'ls > ../py/list_of_files.txt'
os.system(cmd)

file_path_py = '../py/list_of_files.txt'
access_mode = 'r+'

l = open(file_path_py, access_mode)
file_names = l.readlines()

i = 0

while (i < 256):
   next_file = file_names[i]
   next_file = next_file.strip('\n')
   print next_file
   f = open(next_file, access_mode)
   contents = f.read()
   f.close()
   # find anchor code
   f1 = re.search(b'\xD3\xF8\x25\xFF\x00\x00\x00\x50\x68', contents)
   anchor = f1.start()
   print anchor
   # replace put instructions with NOPs
   f1_insert_start = anchor - 55
   f1_insert_len = len(replacement_hex_1) / 2
   f1_insert_end = f1_insert_start + f1_insert_len
   mod = contents[0:f1_insert_start]
   shellcode1 = replacement_hex_1.decode('hex')
   rest = contents[f1_insert_end:]
   contents = mod + shellcode1 + rest
   # replace scanf instructions with NOPs
   f2_insert_start = anchor - 38
   f2_insert_len = len(replacement_hex_2) / 2
   f2_insert_end = f2_insert_start + f2_insert_len
   mod = contents[0:f2_insert_start]
   shellcode2 = replacement_hex_2.decode('hex')
   rest = contents[f2_insert_end:]
   contents = mod + shellcode2 + rest
   # replace cmp instructions with NOPs
   f3_insert_start = anchor - 14
   f3_insert_len = len(replacement_hex_3) / 2
   f3_insert_end = f3_insert_start + f3_insert_len
   mod = contents[0:f3_insert_start]
   shellcode3 = replacement_hex_3.decode('hex')
   rest = contents[f3_insert_end:]
   contents = mod + shellcode3 + rest
   # replace jnz instruction with NOPs
   f4_insert_start = anchor - 8
   f4_insert_len = len(replacement_hex_4) / 2
   f4_insert_end = f4_insert_start + f4_insert_len
   mod = contents[0:f4_insert_start]
   shellcode4 = replacement_hex_4.decode('hex')
   rest = contents[f4_insert_end:]
   contents = mod + shellcode4 + rest
   # write to file
   file_path_mod = '../mod/'
   next_file_name = file_path_mod + 'mod_' + next_file
   g = open(next_file_name, 'wb')
   g.write(contents)
   g.close()
   i = i + 1




