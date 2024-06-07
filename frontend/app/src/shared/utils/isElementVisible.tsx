import { type RefObject, useEffect, useState } from 'react';

export const useIsVisible = (ref: RefObject<HTMLDivElement> | null) => {
  const [isIntersecting, setIntersecting] = useState(false);

  console.log('ref in useIsVisible', ref);

  useEffect(() => {
    if (ref?.current) {
      const observer = new IntersectionObserver(([entry]) => {
        setIntersecting(entry.isIntersecting);
      });

      observer.observe(ref.current);
      return () => {
        observer.disconnect();
      };
    }
  }, [ref]);

  return isIntersecting;
};
