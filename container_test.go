package goroaring

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func makeContainer(ss []short) Container {
	c := NewArrayContainer()
	//	c := NewBitmapContainer()
	for _, s := range ss {
		c.Add(s)
	}
	return c
}

func checkContent(c Container, s []short) bool {
	si := c.GetShortIterator()
	ctr := 0
	fail := false
	for si.HasNext() {
		if ctr == len(s) {
			log.Println("HERE")
			fail = true
			break
		}
		i := si.Next()
		if i != s[ctr] {

			log.Println("THERE", i, s[ctr])
			fail = true
			break
		}
		ctr++
	}
	if ctr != len(s) {
		log.Println("LAST")
		fail = true
	}
	if fail {
		log.Println("fail, found ")
		si = c.GetShortIterator()
		z := 0
		for si.HasNext() {
			//			log.Println(" ", z, si.Next(), s[z])
			si.Next()
			z++
		}
		log.Println(z, len(s))
	}

	return !fail
}

func TestRoaringContainer(t *testing.T) {
	Convey("NumberOfTrailingZeros", t, func() {
		x := int64(0)
		o := NumberOfTrailingZeros(x)
		So(o, ShouldEqual, 64)
		x = 1 << 3
		o = NumberOfTrailingZeros(x)
		So(o, ShouldEqual, 3)
	})
	Convey("ArrayShortIterator", t, func() {
		content := []short{1, 3, 5, 7, 9}
		c := makeContainer(content)
		si := c.GetShortIterator()
		i := 0
		for si.HasNext() {
			si.Next()
			i++
		}

		So(i, ShouldEqual, 5)
	})

	Convey("BinarySearch", t, func() {
		content := []short{1, 3, 5, 7, 9}
		res := binarySearch(content, 5)
		So(res, ShouldEqual, 2)
		res = binarySearch(content, 4)
		So(res, ShouldBeLessThan, 0)
	})
	Convey("bitmapcontainer", t, func() {
		content := []short{1, 3, 5, 7, 9}
		a := NewArrayContainer()
		b := NewBitmapContainer()
		for _, v := range content {
			a.Add(v)
			b.Add(v)
		}
		c := a.ToBitmapContainer()

		So(a.GetCardinality(), ShouldEqual, b.GetCardinality())
		So(c.GetCardinality(), ShouldEqual, b.GetCardinality())

	})
	Convey("inottest0", t, func() {
		content := []short{9}
		c := makeContainer(content)
		c = c.Inot(0, 10)
		si := c.GetShortIterator()
		i := 0
		for si.HasNext() {
			si.Next()
			i++
		}
		So(i, ShouldEqual, 10)
	})

	Convey("inotTest1", t, func() {
		// Array container, range is complete
		content := []short{1, 3, 5, 7, 9}
		//content := []short{1}
		edge := 1 << 13
		//		edge := 30
		c := makeContainer(content)
		c = c.Inot(0, edge)
		size := (edge + 0) - len(content)
		s := make([]short, size)
		pos := 0
		for i := short(0); i < short(edge+0); i++ {
			if binarySearch(content, i) < 0 {
				s[pos] = i
				pos++
			}
		}
		So(checkContent(c, s), ShouldEqual, true)
	})

	/*
	   @Test
	   public void inotTest10() {
	           System.out.println("inotTest10");
	           // Array container, inverting a range past any set bit
	           final short[] content = new short[3];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           final Container c = makeContainer(content);
	           final Container c1 = c.inot(65190, 65200);
	           assertTrue(c1 instanceof ArrayContainer);
	           assertEquals(14, c1.getCardinality());
	           assertTrue(checkContent(c1, new short[] { 0, 2, 4,
	                   (short) 65190, (short) 65191, (short) 65192,
	                   (short) 65193, (short) 65194, (short) 65195,
	                   (short) 65196, (short) 65197, (short) 65198,
	                   (short) 65199, (short) 65200 }));
	   }

	   @Test
	   public void inotTest2() {
	           // Array and then Bitmap container, range is complete
	           final short[] content = { 1, 3, 5, 7, 9 };
	           Container c = makeContainer(content);
	           c = c.inot(0, 65535);
	           c = c.inot(0, 65535);
	           assertTrue(checkContent(c, content));
	   }

	   @Test
	   public void inotTest3() {
	           // Bitmap to bitmap, full range

	           Container c = new ArrayContainer();
	           for (int i = 0; i < 65536; i += 2)
	                   c = c.add((short) i);

	           c = c.inot(0, 65535);
	           assertTrue(c.contains((short) 3) && !c.contains((short) 4));
	           assertEquals(32768, c.getCardinality());
	           c = c.inot(0, 65535);
	           for (int i = 0; i < 65536; i += 2)
	                   assertTrue(c.contains((short) i)
	                           && !c.contains((short) (i + 1)));
	   }

	   @Test
	   public void inotTest4() {
	           // Array container, range is partial, result stays array
	           final short[] content = { 1, 3, 5, 7, 9 };
	           Container c = makeContainer(content);
	           c = c.inot(4, 999);
	           assertTrue(c instanceof ArrayContainer);
	           assertEquals(999 - 4 + 1 - 3 + 2, c.getCardinality());
	           c = c.inot(4, 999); // back
	           assertTrue(checkContent(c, content));
	   }

	   @Test
	   public void inotTest5() {
	           System.out.println("inotTest5");
	           // Bitmap container, range is partial, result stays bitmap
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           Container c = makeContainer(content);
	           c = c.inot(4, 999);
	           assertTrue(c instanceof BitmapContainer);
	           assertEquals(31773, c.getCardinality());
	           c = c.inot(4, 999); // back, as a bitmap
	           assertTrue(c instanceof BitmapContainer);
	           assertTrue(checkContent(c, content));

	   }

	   @Test
	   public void inotTest6() {
	           System.out.println("inotTest6");
	           // Bitmap container, range is partial and in one word, result
	           // stays bitmap
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           Container c = makeContainer(content);
	           c = c.inot(4, 8);
	           assertTrue(c instanceof BitmapContainer);
	           assertEquals(32762, c.getCardinality());
	           c = c.inot(4, 8); // back, as a bitmap
	           assertTrue(c instanceof BitmapContainer);
	           assertTrue(checkContent(c, content));
	   }

	   @Test
	   public void inotTest7() {
	           System.out.println("inotTest7");
	           // Bitmap container, range is partial, result flips to array
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           Container c = makeContainer(content);
	           c = c.inot(5, 31000);
	           if(c.getCardinality() <= ArrayContainer.DEFAULTMAXSIZE)
	                   assertTrue(c instanceof ArrayContainer);
	           else
	                   assertTrue(c instanceof BitmapContainer);
	           assertEquals(1773, c.getCardinality());
	           c = c.inot(5, 31000); // back, as a bitmap
	           if(c.getCardinality() <= ArrayContainer.DEFAULTMAXSIZE)
	                   assertTrue(c instanceof ArrayContainer);
	           else
	                   assertTrue(c instanceof BitmapContainer);
	           assertTrue(checkContent(c, content));
	   }

	   // case requiring contraction of ArrayContainer.
	   @Test
	   public void inotTest8() {
	           System.out.println("inotTest8");
	           // Array container
	           final short[] content = new short[21];
	           for (int i = 0; i < 18; ++i)
	                   content[i] = (short) i;
	           content[18] = 21;
	           content[19] = 22;
	           content[20] = 23;

	           Container c = makeContainer(content);
	           c = c.inot(5, 21);
	           assertTrue(c instanceof ArrayContainer);

	           assertEquals(10, c.getCardinality());
	           c = c.inot(5, 21); // back, as a bitmap
	           assertTrue(c instanceof ArrayContainer);
	           assertTrue(checkContent(c, content));
	   }

	   // mostly same tests, except for not. (check original unaffected)
	   @Test
	   public void notTest1() {
	           // Array container, range is complete
	           final short[] content = { 1, 3, 5, 7, 9 };
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(0, 65535);
	           final short[] s = new short[65536 - content.length];
	           int pos = 0;
	           for (int i = 0; i < 65536; ++i)
	                   if (Arrays.binarySearch(content, (short) i) < 0)
	                           s[pos++] = (short) i;
	           assertTrue(checkContent(c1, s));
	           assertTrue(checkContent(c, content));
	   }

	   @Test
	   public void notTest10() {
	           System.out.println("notTest10");
	           // Array container, inverting a range past any set bit
	           // attempting to recreate a bug (but bug required extra space
	           // in the array with just the right junk in it.
	           final short[] content = new short[40];
	           for (int i = 244; i <= 283; ++i)
	                   content[i - 244] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(51413, 51470);
	           assertTrue(c1 instanceof ArrayContainer);
	           assertEquals(40 + 58, c1.getCardinality());
	           final short[] rightAns = new short[98];
	           for (int i = 244; i <= 283; ++i)
	                   rightAns[i - 244] = (short) i;
	           for (int i = 51413; i <= 51470; ++i)
	                   rightAns[i - 51413 + 40] = (short) i;

	           assertTrue(checkContent(c1, rightAns));
	   }

	   @Test
	   public void notTest11() {
	           System.out.println("notTest11");
	           // Array container, inverting a range before any set bit
	           // attempting to recreate a bug (but required extra space
	           // in the array with the right junk in it.
	           final short[] content = new short[40];
	           for (int i = 244; i <= 283; ++i)
	                   content[i - 244] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(1, 58);
	           assertTrue(c1 instanceof ArrayContainer);
	           assertEquals(40 + 58, c1.getCardinality());
	           final short[] rightAns = new short[98];
	           for (int i = 1; i <= 58; ++i)
	                   rightAns[i - 1] = (short) i;
	           for (int i = 244; i <= 283; ++i)
	                   rightAns[i - 244 + 58] = (short) i;

	           assertTrue(checkContent(c1, rightAns));
	   }

	   @Test
	   public void notTest2() {
	           // Array and then Bitmap container, range is complete
	           final short[] content = { 1, 3, 5, 7, 9 };
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(0, 65535);
	           final Container c2 = c1.not(0, 65535);
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void notTest3() {
	           // Bitmap to bitmap, full range

	           Container c = new ArrayContainer();
	           for (int i = 0; i < 65536; i += 2)
	                   c = c.add((short) i);

	           final Container c1 = c.not(0, 65535);
	           assertTrue(c1.contains((short) 3) && !c1.contains((short) 4));
	           assertEquals(32768, c1.getCardinality());
	           final Container c2 = c1.not(0, 65535);
	           for (int i = 0; i < 65536; i += 2)
	                   assertTrue(c2.contains((short) i)
	                           && !c2.contains((short) (i + 1)));
	   }

	   @Test
	   public void notTest4() {
	           System.out.println("notTest4");
	           // Array container, range is partial, result stays array
	           final short[] content = { 1, 3, 5, 7, 9 };
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(4, 999);
	           assertTrue(c1 instanceof ArrayContainer);
	           assertEquals(999 - 4 + 1 - 3 + 2, c1.getCardinality());
	           final Container c2 = c1.not(4, 999); // back
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void notTest5() {
	           System.out.println("notTest5");
	           // Bitmap container, range is partial, result stays bitmap
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(4, 999);
	           assertTrue(c1 instanceof BitmapContainer);
	           assertEquals(31773, c1.getCardinality());
	           final Container c2 = c1.not(4, 999); // back, as a bitmap
	           assertTrue(c2 instanceof BitmapContainer);
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void notTest6() {
	           System.out.println("notTest6");
	           // Bitmap container, range is partial and in one word, result
	           // stays bitmap
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(4, 8);
	           assertTrue(c1 instanceof BitmapContainer);
	           assertEquals(32762, c1.getCardinality());
	           final Container c2 = c1.not(4, 8); // back, as a bitmap
	           assertTrue(c2 instanceof BitmapContainer);
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void notTest7() {
	           System.out.println("notTest7");
	           // Bitmap container, range is partial, result flips to array
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(5, 31000);
	           if(c1.getCardinality() <= ArrayContainer.DEFAULTMAXSIZE)
	                   assertTrue(c1 instanceof ArrayContainer);
	           else
	                   assertTrue(c1 instanceof BitmapContainer);
	           assertEquals(1773, c1.getCardinality());
	           final Container c2 = c1.not(5, 31000); // back, as a bitmap
	           if(c2.getCardinality() <= ArrayContainer.DEFAULTMAXSIZE)
	                   assertTrue(c2 instanceof ArrayContainer);
	           else
	                   assertTrue(c2 instanceof BitmapContainer);
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void notTest8() {
	           System.out.println("notTest8");
	           // Bitmap container, range is partial on the lower end
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(4, 65535);
	           assertTrue(c1 instanceof BitmapContainer);
	           assertEquals(32773, c1.getCardinality());
	           final Container c2 = c1.not(4, 65535); // back, as a bitmap
	           assertTrue(c2 instanceof BitmapContainer);
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void notTest9() {
	           System.out.println("notTest9");
	           // Bitmap container, range is partial on the upper end, not
	           // single word
	           final short[] content = new short[32768 - 5];
	           content[0] = 0;
	           content[1] = 2;
	           content[2] = 4;
	           content[3] = 6;
	           content[4] = 8;
	           for (int i = 10; i <= 32767; ++i)
	                   content[i - 10 + 5] = (short) i;
	           final Container c = makeContainer(content);
	           final Container c1 = c.not(0, 65200);
	           assertTrue(c1 instanceof BitmapContainer);
	           assertEquals(32438, c1.getCardinality());
	           final Container c2 = c1.not(0, 65200); // back, as a bitmap
	           assertTrue(c2 instanceof BitmapContainer);
	           assertTrue(checkContent(c2, content));
	   }

	   @Test
	   public void rangeOfOnesTest1() {
	           final Container c = Container.rangeOfOnes(4, 10); // sparse
	           assertTrue(c instanceof ArrayContainer);
	           assertEquals(10 - 4 + 1, c.getCardinality());
	           assertTrue(checkContent(c, new short[] { 4, 5, 6, 7, 8, 9, 10 }));
	   }

	   @Test
	   public void rangeOfOnesTest2() {
	           final Container c = Container.rangeOfOnes(1000, 35000); // dense
	           assertTrue(c instanceof BitmapContainer);
	           assertEquals(35000 - 1000 + 1, c.getCardinality());
	   }

	   @Test
	   public void rangeOfOnesTest2A() {
	           final Container c = Container.rangeOfOnes(1000, 35000); // dense
	           final short s[] = new short[35000 - 1000 + 1];
	           for (int i = 1000; i <= 35000; ++i)
	                   s[i - 1000] = (short) i;
	           assertTrue(checkContent(c, s));
	   }

	   @Test
	   public void rangeOfOnesTest3() {
	           // bdry cases
	           final Container c = Container.rangeOfOnes(1,
	                   ArrayContainer.DEFAULTMAXSIZE);
	           assertTrue(c instanceof ArrayContainer);
	   }

	   @Test
	   public void rangeOfOnesTest4() {
	           final Container c = Container.rangeOfOnes(1,
	                   ArrayContainer.DEFAULTMAXSIZE + 1);
	           assertTrue(c instanceof BitmapContainer);
	   }

	   public static boolean checkContent(Container c, short[] s) {
	           ShortIterator si = c.getShortIterator();
	           int ctr = 0;
	           boolean fail = false;
	           while (si.hasNext()) {
	                   if (ctr == s.length) {
	                           fail = true;
	                           break;
	                   }
	                   if (si.next() != s[ctr]) {
	                           fail = true;
	                           break;
	                   }
	                   ++ctr;
	           }
	           if (ctr != s.length) {
	                   fail = true;
	           }
	           if (fail) {
	                   System.out.print("fail, found ");
	                   si = c.getShortIterator();
	                   while (si.hasNext())
	                           System.out.print(" " + si.next());
	                   System.out.print("\n expected ");
	                   for (final short s1 : s)
	                           System.out.print(" " + s1);
	                   System.out.println();
	           }
	           return !fail;
	   }

	*/

}
